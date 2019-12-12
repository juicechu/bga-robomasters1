package app

/*
#include <stdio.h>

static void callback(unsigned long long e, void* info, unsigned long long tag) {
	printf("Unity bridge callback called!\n");
}
*/
import "C"
import (
	"encoding/binary"
	"fmt"
	"io/ioutil"

	"git.bug-br.org.br/bga/robomasters1/app/internal/dji/unitybridge/wrapper"
	"git.bug-br.org.br/bga/robomasters1/app/internal/pairing"
	"git.bug-br.org.br/bga/robomasters1/app/internal/udp"
	"github.com/google/uuid"
	"github.com/skip2/go-qrcode"
	"github.com/skratchdot/open-golang/open"

	internalqrcode "git.bug-br.org.br/bga/robomasters1/app/internal/qrcode"
)

type App struct {
	id  uint64
	qrc *internalqrcode.QRCode
	pl  *pairing.Listener
	pp  *udp.PortPair
}

func New(countryCode, ssId, password, bssId string) (*App, error) {
	appId, err := generateAppId()
	if err != nil {
		return nil, err
	}

	return NewWithAppID(countryCode, ssId, password, bssId, appId)
}

func NewWithAppID(countryCode, ssId, password, bssId string,
	appId uint64) (*App, error) {
	qrc, err := internalqrcode.NewQRCode(appId, countryCode, ssId, password,
		bssId)
	if err != nil {
		return nil, err
	}

	return &App{
		appId,
		qrc,
		pairing.NewListener(appId),
		udp.NewPortPair(10607, 10609, 1514),
	}, nil
}

func (a *App) Start(textMode bool) error {
	var err error
	if textMode {
		err = a.showTextQRCode()
	} else {
		err = a.showPNGQRCode()
	}
	if err != nil {
		return fmt.Errorf("error showing QR code: %w", err)
	}

	eventChan, err := a.pl.Start()
	if err != nil {
		return fmt.Errorf("error starting pairing listener: %w", err)
	}

	packetChan, err := a.pp.Start()
	if err != nil {
		return fmt.Errorf("error starting packet listener: %w", err)
	}

	// Setup Unity Bridge.
	//
	// TODO(bga): Move this to its own function/method.
	wrapper.Instance().CreateUnityBridge("Robomaster", true);
	eventTypes := []uint64{0, 1, 2, 3, 4, 5, 6, 7, 8, 100, 101, 200, 300,
		301, 302, 303, 304, 305,306, 500}
	for eventType := range eventTypes {
		wrapper.Instance().RegisterEventCallback(eventType << 32, C.callback)
	}
	ok := wrapper.Instance().UnityBridgeInitialize()
	if !ok {
		wrapper.Instance().DestroyUnityBridge()
		return fmt.Errorf("failed initializing unity bridge")
	}

L:
	for {
		select {
		case event, ok := <-eventChan:
			if !ok {
				break L
			}

			// TODO(bga): Do something meaningful.
			fmt.Printf("%#+v\n", *event)
		case packet, ok := <-packetChan:
			if !ok {
				break L
			}

			fmt.Printf("%#+v\n", *packet)
		}
	}

	return nil
}

func (a *App) showTextQRCode() error {
	qrc, err := qrcode.New(a.qrc.EncodedMessage(), qrcode.Medium)
	if err != nil {
		return err
	}

	fmt.Println(qrc.ToString(false))

	return nil
}

func (a *App) showPNGQRCode() error {
	pngData, err := qrcode.Encode(a.qrc.EncodedMessage(), qrcode.Medium,
		256)
	if err != nil {
		return err
	}

	f, err := ioutil.TempFile("", "qrcode-*.png")
	if err != nil {
		return err
	}

	fileName := f.Name()

	_, err = f.Write(pngData)
	if err != nil {
		f.Close()
		return err
	}

	f.Close()

	return open.Run(fileName)
}

func generateAppId() (uint64, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return 0, err
	}

	// Create an app ID out of the first 8 bytes of the UUID.
	return binary.LittleEndian.Uint64(id[0:9]), nil
}
