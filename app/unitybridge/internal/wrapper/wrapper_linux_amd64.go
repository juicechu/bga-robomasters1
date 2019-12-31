package wrapper

import (
	"encoding/binary"
	"git.bug-br.org.br/bga/robomasters1/app/unitybridge/internal/wrapper/winebridge"
	"io"
	"log"
	"os"
)

type Linux struct {
	*generic

	readPipe  io.Reader
	writePipe io.Writer

	wineBridge *winebridge.WineBridge
}

func Instance() Wrapper {
	once.Do(func() {
		localReadPipe, remoteWritePipe, err := os.Pipe()
		if err != nil {
			panic(err)
		}

		remoteReadPipe, localWritePipe, err := os.Pipe()
		if err != nil {
			panic(err)
		}

		wineBridge, err := winebridge.New("winewrapper.exe",
			remoteReadPipe, remoteWritePipe)

		err = wineBridge.Start()
		if err != nil {
			panic(err)
		}

		l := &Linux{
			readPipe:   localReadPipe,
			writePipe:  localWritePipe,
			wineBridge: wineBridge,
		}

		l.generic = newGeneric(l)

		go l.readLoop()

		instance = l
	})

	return instance
}

func (l *Linux) CreateUnityBridge(name string, debuggable bool) {
	// size (4 bytes) + function (1 byte) + debuggable (1 byte) + len(name)
	buffer := make([]byte, 4+1+1+len(name))
	binary.LittleEndian.PutUint32(buffer, uint32(len(buffer)))
	buffer[4] = FuncCreateUnityBridge
	if debuggable {
		buffer[5] = 1
	} else {
		buffer[5] = 0
	}
	copy(buffer[6:], name)

	_, err := l.writePipe.Write(buffer)
	if err != nil {
		panic(err)
	}
}

func (l *Linux) DestroyUnityBridge() {
	// size (4 bytes) + function (1 byte)
	buffer := make([]byte, 4+1)
	binary.LittleEndian.PutUint32(buffer, uint32(len(buffer)))
	buffer[4] = FuncDestroyUnityBridge

	_, err := l.writePipe.Write(buffer)
	if err != nil {
		panic(err)
	}
}

func (l *Linux) UnityBridgeInitialize() bool {
	// size (4 bytes) + function (1 byte)
	buffer := make([]byte, 4+1)
	binary.LittleEndian.PutUint32(buffer, uint32(len(buffer)))
	buffer[4] = FuncUnityBridgeInitialize

	_, err := l.writePipe.Write(buffer)
	if err != nil {
		panic(err)
	}

	return true
}

func (l *Linux) UnityBridgeUninitialize() {
	// size (4 bytes) + function (1 byte)
	buffer := make([]byte, 4+1)
	binary.LittleEndian.PutUint32(buffer, uint32(len(buffer)))
	buffer[4] = FuncUnityBridgeUninitialize

	_, err := l.writePipe.Write(buffer)
	if err != nil {
		panic(err)
	}
}

func (l *Linux) UnitySendEvent(eventCode uint64, info []byte, tag uint64) {
	// size (4 bytes) + function (1 byte) + eventCode (8 bytes) +
	// tag (8 bytes) + len(info)
	buffer := make([]byte, 4+1+8+8+len(info))
	binary.LittleEndian.PutUint32(buffer, uint32(len(buffer)))
	buffer[4] = FuncUnitySendEvent
	binary.LittleEndian.PutUint64(buffer[5:], eventCode)
	binary.LittleEndian.PutUint64(buffer[13:], tag)
	copy(buffer[21:], info)

	_, err := l.writePipe.Write(buffer)
	if err != nil {
		panic(err)
	}
}

func (l *Linux) UnitySendEventWithNumber(eventCode, info, tag uint64) {
	// size (4 bytes) + function (1 byte) + eventCode (8 bytes) +
	// tag (8 bytes) + info (8 bytes)
	buffer := make([]byte, 4+1+8+8+8)
	binary.LittleEndian.PutUint32(buffer, uint32(len(buffer)))
	buffer[4] = FuncUnitySendEventWithNumber
	binary.LittleEndian.PutUint64(buffer[5:], eventCode)
	binary.LittleEndian.PutUint64(buffer[13:], tag)
	binary.LittleEndian.PutUint64(buffer[21:], info)

	_, err := l.writePipe.Write(buffer)
	if err != nil {
		panic(err)
	}
}

func (l *Linux) UnitySendEventWithString(eventCode uint64, info string,
	tag uint64) {
	// size (4 bytes) + function (1 byte) + eventCode (8 bytes) +
	// tag (8 bytes) + len(info)
	buffer := make([]byte, 4+1+8+8+len(info))
	binary.LittleEndian.PutUint32(buffer, uint32(len(buffer)))
	buffer[4] = FuncUnitySendEventWithString
	binary.LittleEndian.PutUint64(buffer[5:], eventCode)
	binary.LittleEndian.PutUint64(buffer[13:], tag)
	copy(buffer[21:], info)

	_, err := l.writePipe.Write(buffer)
	if err != nil {
		panic(err)
	}
}

func (l *Linux) unitySetEventCallback(eventCode uint64, add bool) {
	// size (4 bytes) + function (1 byte) + add (1 byte) +
	// eventCode (8 bytes)
	buffer := make([]byte, 4+1+1+8)
	binary.LittleEndian.PutUint32(buffer, uint32(len(buffer)))
	buffer[4] = FuncUnitySetEventCallback
	if add {
		buffer[5] = 1
	} else {
		buffer[5] = 0
	}
	binary.LittleEndian.PutUint64(buffer[6:], eventCode)

	_, err := l.writePipe.Write(buffer)
	if err != nil {
		panic(err)
	}
}

func (l *Linux) readLoop() {
	readBuffer := make([]byte, 10000)
	lengthBuffer := readBuffer[:4]
	for {
		_, err := io.ReadFull(l.readPipe, lengthBuffer)
		if err != nil {
			panic(err)
		}

		length := binary.LittleEndian.Uint32(lengthBuffer)

		length -= 4

		if len(readBuffer) < int(length) {
			readBuffer = make([]byte, length)
		}

		sizedReadBuffer := readBuffer[:length]

		_, err = io.ReadFull(l.readPipe, sizedReadBuffer)
		if err != nil {
			panic(err)
		}

		eventCode := binary.LittleEndian.Uint64(sizedReadBuffer)
		tag := binary.LittleEndian.Uint64(sizedReadBuffer[8:])
		info := sizedReadBuffer[16:]

		callback := l.Callback(eventCode)
		if callback == nil {
			log.Printf("No callback for event code %d.\n",
				eventCode)
		}

		callback(eventCode, info, tag)
	}
}

func getFd(file *os.File) uintptr {
	rawConn, err := file.SyscallConn()
	if err != nil {
		panic(err)
	}

	var fileFd uintptr
	err = rawConn.Control(func(fd uintptr) {
		fileFd = fd
	})
	if err != nil {
		panic(err)
	}

	return fileFd
}
