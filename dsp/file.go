package dsp

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"git.bug-br.org.br/bga/robomasters1/dsp/internal"
)

var (
	// Extacted from DJI's RoboMaster S1 app.
	dspKey = []byte("TRoP4GWuc30k6WUp")
	dspIv  = []byte("bP3crVEO6wABzOc0")
)

// File is the representation of a RoboMaster S1 program file (.dsp). It can be
// used to create, read, modify or write them.
type File struct {
	dji internal.Dji
}

// New creates a new File instance with the given creator, title and pythonCode.
// Returns a pointer to a File instance and a nil error on success or nil and a
// non-nil error on failure.
func New(creator, title, pythonCode string) (*File, error) {
	trimmedCreator := strings.TrimSpace(creator)
	if len(trimmedCreator) == 0 {
		return nil, fmt.Errorf("creator can not be empty")
	}

	trimmedTitle := strings.TrimSpace(title)
	if len(trimmedTitle) == 0 {
		return nil, fmt.Errorf("title can not be empty")
	}

	trimmedPythonCode := strings.TrimSpace(pythonCode)
	if len(trimmedCreator) == 0 {
		return nil, fmt.Errorf("pythonCode can not be empty")
	}

	now := time.Now()

	f := &File{
		internal.Dji{
			Attribute: internal.Attribute{
				Creator:      trimmedCreator,
				CreationDate: now.Format("2006/01/02"),
				ModifyTime: now.Format(
					"01/02/2006 15:04:05 MST"),
				FirmwareVersionDependency: "00.00.0000",
				Title:                     trimmedTitle,
				Guid:                      computeGuid(),
				CodeType:                  "python",
				AppMinVersion:             "",
				AppMaxVersion:             "",
				Sign:                      "",
				// To be signed after construction.
			},
			Code: internal.Code{
				PythonCode: internal.Cdata{
					Cdata: trimmedPythonCode,
				},
				ScratchDescription: internal.Cdata{
					Cdata: "",
				},
			},
		},
	}

	f.computeSignature()

	return f, nil
}

// Load loads a RoboMaster S1 program file (.dsp) from disk. Returns a pointer
// to a File instance and a nil error on success or nil and a non-nil error on
// failure.
func Load(fileName string) (*File, error) {
	xmlData, err := decodeDsp(fileName)
	if err != nil {
		return nil, err
	}

	var f File
	err = xml.Unmarshal(xmlData, &f.dji)

	return &f, nil
}

// Save serializes and saves the File instance to disk as an encrypted
// RoboMaster S1 program file (.dsp). Returns a nil error on success or a
// non-nil error on failure.
func (f *File) Save(fileName string) error {
	fd, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer fd.Close()

	xmlData, err := xml.Marshal(f.dji)
	if err != nil {
		return err
	}

	dspData, err := encodeDsp(xmlData)
	if err != nil {
		return err
	}

	_, err = fd.Write(dspData)
	if err != nil {
		return err
	}

	return nil
}

func (f *File) computeSignature() {
	// TODO(bga): Actually compute signature.
	f.dji.Attribute.Sign = "signature"
}

func computeGuid() string {
	// TODO(bga): Do something reasonable.
	return "guid"
}

func decodeDsp(fileName string) ([]byte, error) {
	fd, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	decoder := base64.NewDecoder(base64.StdEncoding, fd)

	cipherText, err := ioutil.ReadAll(decoder)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(dspKey)
	if err != nil {
		return nil, err
	}

	cbc := cipher.NewCBCDecrypter(block, dspIv)
	cbc.CryptBlocks(cipherText, cipherText)

	return cipherText, nil
}

func encodeDsp(plainText []byte) ([]byte, error) {
	extraBytes := len(plainText) % aes.BlockSize
	if extraBytes != 0 {
		paddingBytes := aes.BlockSize - extraBytes
		padding := []byte(bytes.Repeat([]byte{' '}, paddingBytes))
		plainText = append(plainText, padding...)
	}

	block, err := aes.NewCipher(dspKey)
	if err != nil {
		return nil, err
	}

	cbc := cipher.NewCBCEncrypter(block, dspIv)

	cipherText := make([]byte, len(plainText))
	cbc.CryptBlocks(cipherText, plainText)

	base64Enc := base64.StdEncoding

	base64Text := make([]byte, base64Enc.EncodedLen(len(cipherText)))
	base64Enc.Encode(base64Text, cipherText)

	return base64Text, nil
}