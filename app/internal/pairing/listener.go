package pairing

import (
	"encoding/binary"
	"fmt"
	"net"
	"sync"
)

const (
	listenerServerPort = 45678
	listenerClientPort = 56789
	maxBufferSize      = 256
)

type Listener struct {
	appId uint64

	m          sync.Mutex
	packetConn net.PacketConn
	eventChan  chan *Event
	quitChan   chan struct{}
	clientMap  map[string]bool

	wg sync.WaitGroup
}

func NewListener(appId uint64) *Listener {
	return &Listener{
		appId,
		sync.Mutex{},
		nil,
		nil,
		nil,
		make(map[string]bool),
		sync.WaitGroup{},
	}
}

func (l *Listener) Start() (<-chan *Event, error) {
	l.m.Lock()
	defer l.m.Unlock()

	if l.packetConn != nil {
		return nil, fmt.Errorf("already started")
	}

	packetConn, err := net.ListenPacket("udp", fmt.Sprintf(
		":%d", listenerServerPort))
	if err != nil {
		return nil, fmt.Errorf("error starting listener: %w", err)
	}

	packetConn.(*net.UDPConn).SetReadBuffer(maxBufferSize)
	packetConn.(*net.UDPConn).SetWriteBuffer(maxBufferSize)

	l.packetConn = packetConn

	l.eventChan = make(chan *Event)
	l.quitChan = make(chan struct{})

	l.wg.Add(1)
	go l.loop()

	return l.eventChan, nil
}

func (l *Listener) Stop() error {
	l.m.Lock()
	defer l.m.Unlock()

	if l.packetConn == nil {
		return fmt.Errorf("not started")
	}

	close(l.quitChan)
	l.packetConn.Close()

	l.packetConn = nil

	l.eventChan = nil
	l.quitChan = nil

	l.wg.Wait()

	return nil
}

func (l *Listener) sendACK(ip net.IP) error {
	buffer := make([]byte, 8)
	binary.LittleEndian.PutUint64(buffer, l.appId)

	clientAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", ip.String(),
		listenerClientPort))
	if err != nil {
		return fmt.Errorf("error resolving client address: %w", err)
	}

	_, err = l.packetConn.WriteTo(buffer, clientAddr)
	if err != nil {
		return fmt.Errorf("error sending client ack: %w", err)
	}

	return nil
}

func (l *Listener) maybeGenerateEvent(addr net.Addr, buffer []byte) *Event {
	bm, err := ParseBroadcastMessageData(buffer)
	if err != nil {
		return nil
	}

	l.m.Lock()
	defer l.m.Unlock()

	ip := addr.(*net.UDPAddr).IP
	if l.clientMap[ip.String()] {
		if bm.AppId() != l.appId {
			l.clientMap[ip.String()] = false
			return NewEvent(EventRemove, bm.SourceIp(), bm.SourceMac())
		}
	} else {
		if bm.IsPairing() && bm.AppId() == l.appId {
			err = l.sendACK(ip)
			if err != nil {
				return nil
			}

			l.clientMap[ip.String()] = true
			return NewEvent(EventAdd, bm.SourceIp(), bm.SourceMac())
		}
	}

	return nil
}

func (l *Listener) loop() {
	fullStop := false

	buffer := make([]byte, maxBufferSize)
L:
	for {
		n, addr, err := l.packetConn.ReadFrom(buffer)
		if err != nil {
			fullStop = true
			break L
		}

		if event := l.maybeGenerateEvent(addr, buffer[:n]); event != nil {
			select {
			case <-l.quitChan:
				break L

			case l.eventChan <- event:
				// Do nothing.
			}
		} else {
			select {
			case <-l.quitChan:
				break L
			default:
				// Do nothing.
			}
		}
	}

	close(l.eventChan)
	l.wg.Done()

	if fullStop {
		l.Stop()
	}
}
