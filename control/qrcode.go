package control

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"io/ioutil"

	"github.com/skratchdot/open-golang/open"

	qrcode "github.com/skip2/go-qrcode"
)

type QRCode struct {
	message string
}

func NewQRCode(ssid, password string, appId uint64) (*QRCode, error) {
	message, err := encodeQRCodeMessage(ssid, password, appId)
	if err != nil {
		return nil, err
	}

	return &QRCode{
		message,
	}, nil
}

func (q *QRCode) Show() error {
	pngData, err := qrcode.Encode(q.message, qrcode.Medium, 256)
	if err != nil {
		return err
	}

	f, err := ioutil.TempFile("", "qrcode-*.png")
	if err != nil {
		return err
	}

	fileName := f.Name()

	_, err = f.Write(pngData)
	f.Close()
	if err != nil {
		return err
	}

	return open.Run(fileName)
}

func encodeQRCodeMessage(ssid, password string, appId uint64) (string, error) {
	var b bytes.Buffer

	bytesSsid := []byte(ssid)
	bytesPassword := []byte(password)

	hasBssid := uint16(0)

	metadata := (hasBssid << 11) | (uint16(len(bytesPassword)) << 6) |
		uint16(len(bytesSsid))

	data1 := make([]byte, 2)
	binary.LittleEndian.PutUint16(data1, metadata)

	b.Write(data1)

	data2 := make([]byte, 8)
	binary.LittleEndian.PutUint64(data2, appId)

	b.Write(data2)

	b.Write([]byte("US"))

	b.Write([]byte(ssid))
	b.Write([]byte(password))

	// If we had a BSSID, we would append it here.

	data := b.Bytes()
	inPlaceEncodeDecode(data)

	return base64.StdEncoding.EncodeToString(data), nil
}
