package app

import (
	"encoding/binary"
	"io/ioutil"

	"github.com/google/uuid"
	"github.com/skip2/go-qrcode"
	"github.com/skratchdot/open-golang/open"
)

type App struct {
	id  uint64
	qrc *QRCode
}

func New(countryCode, ssId, password, bssId string) (*App, error) {
	appId, err := generateAppId()
	if err != nil {
		return nil, err
	}

	qrc, err := NewQRCode(appId, countryCode, ssId, password, bssId)
	if err != nil {
		return nil, err
	}

	return &App{
		appId,
		qrc,
	}, nil
}

func (a *App) ShowQRCode() error {
	pngData, err := qrcode.Encode(a.qrc.EncodedMessage(), qrcode.High, 512)
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
