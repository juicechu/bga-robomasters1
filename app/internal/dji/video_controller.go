package dji

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"sync"

	"git.bug-br.org.br/bga/robomasters1/app/internal/dji/unity"
	"git.bug-br.org.br/bga/robomasters1/app/internal/dji/unity/bridge"
	"git.bug-br.org.br/bga/robomasters1/app/internal/rgb"
)

type VideoController struct {
	eventHandlerIndexes []int

	once sync.Once
}

func NewVideoController() (*VideoController, error) {
	b := bridge.Instance()

	vc := &VideoController{}

	eventHandlerIndexes := make([]int, 3)

	var err error
	eventHandlerIndexes[0], err = b.AddEventHandler(
		unity.EventTypeGetNativeTexture, vc)
	if err != nil {
		return nil, err
	}
	eventHandlerIndexes[1], err = b.AddEventHandler(
		unity.EventTypeVideoTransferSpeed, vc)
	if err != nil {
		return nil, err
	}
	eventHandlerIndexes[2], err = b.AddEventHandler(
		unity.EventTypeVideoDataRecv, vc)
	if err != nil {
		return nil, err
	}

	vc.eventHandlerIndexes = eventHandlerIndexes

	return vc, nil
}

func (v *VideoController) StartVideo() {
	ub := bridge.Instance()

	ub.SendEvent(unity.NewEvent(unity.EventTypeStartVideo))
}

func (v *VideoController) StopVideo() {
	ub := bridge.Instance()

	ub.SendEvent(unity.NewEvent(unity.EventTypeStopVideo))
}

func (v *VideoController) HandleEvent(event *unity.Event, info []byte,
	tag uint64) {
	if event.Type() == unity.EventTypeVideoDataRecv {
		v.once.Do(func() {
			rgbImage := &rgb.Image{
				Pix:    info,
				Stride: 3 * 1280,
				Rect:   image.Rect(0, 0, 1280, 720),
			}

			f, err := os.Create("image.png")
			if err != nil {
				panic(err)
			}

			if err := png.Encode(f, rgbImage); err != nil {
				f.Close()
				panic(err)
			}

			if err := f.Close(); err != nil {
				panic(err)
			}
		})
	}
}
