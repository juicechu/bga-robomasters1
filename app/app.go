package app

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"

	"git.bug-br.org.br/bga/robomasters1/app/internal/dji/unitybridge"
	"git.bug-br.org.br/bga/robomasters1/app/internal/pairing"
	"github.com/google/uuid"
	"github.com/skip2/go-qrcode"
	"github.com/skratchdot/open-golang/open"

	internalqrcode "git.bug-br.org.br/bga/robomasters1/app/internal/qrcode"
)

type App struct {
	id  uint64
	qrc *internalqrcode.QRCode
	pl  *pairing.Listener
	ub  *unitybridge.UnityBridge
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
		nil,
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

	// Setup Unity Bridge.
	ub, err := unitybridge.New("Robomaster", true)
	if err != nil {
		return err
	}

	a.ub = ub

	// Start listening to AirlinkConnection events.
	err = a.ub.SendEvent(uint64(4)<<32 | uint64(117440513))
	if err != nil {
		panic(err)
	}

	// Reset connection to defaults.
	err = a.ub.SendEvent((uint64(100)<<32)|uint64(2), "192.168.2.1")
	if err != nil {
		panic(err)
	}
	err = a.ub.SendEvent((uint64(100)<<32)|uint64(3), uint64(10607),
		uint64(0))
	if err != nil {
		panic(err)
	}
	err = a.ub.SendEvent((uint64(100) << 32) | uint64(0))
	if err != nil {
		panic(err)
	}

	eventChan, err := a.pl.Start()
	if err != nil {
		return fmt.Errorf("error starting pairing listener: %w", err)
	}

L:
	for {
		select {
		case event, ok := <-eventChan:
			if !ok {
				break L
			}

			if event.Type() == pairing.EventAdd {
				err = a.ub.SendEvent(
					(uint64(100) << 32) | uint64(1))
				if err != nil {
					panic(err)
				}
				err = a.ub.SendEvent((uint64(100)<<32)|
					uint64(2), event.IP().String())
				if err != nil {
					panic(err)
				}
				err = a.ub.SendEvent((uint64(100)<<32)|
					uint64(3), uint64(10607))
				if err != nil {
					panic(err)
				}
				err = a.ub.SendEvent((uint64(100) << 32) |
					uint64(0))
				if err != nil {
					panic(err)
				}
			}

			fmt.Printf("%#+v\n", event)
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
