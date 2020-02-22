package main

import (
	"image"
	"sync"
	"unsafe"

	"git.bug-br.org.br/bga/robomasters1/app/video"
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

type VideoSystem struct {
	videoEntity      *VideoEntity
	elapsed          float32
	frameCh          chan *image.NRGBA
	dataHandlerIndex int
	wg               *sync.WaitGroup
}

func (s *VideoSystem) New(w *ecs.World) {
	rect := image.Rect(0, 0, 1280, 720)
	s.videoEntity = &VideoEntity{BasicEntity: ecs.NewBasic()}
	s.videoEntity.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{X: 0, Y: 0},
		Width:    1280,
		Height:   720,
	}

	img := image.NewNRGBA(rect)

	obj := common.NewImageObject(img)

	s.videoEntity.RenderComponent = common.RenderComponent{
		Drawable: common.NewTextureSingle(obj),
	}

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&s.videoEntity.BasicEntity,
				&s.videoEntity.RenderComponent,
				&s.videoEntity.SpaceComponent)
		}
	}

	s.frameCh = make(chan *image.NRGBA, 30)

	v, err := video.New()
	if err != nil {
		panic(err)
	}

	index, err := v.AddDataHandler(s.DataHandler)
	if err != nil {
		panic(err)
	}

	s.dataHandlerIndex = index

	v.StartVideo()
}

func (s *VideoSystem) Add() {}

func (s *VideoSystem) Remove(basic ecs.BasicEntity) {}

func (s *VideoSystem) Update(dt float32) {
	select {
	case img := <-s.frameCh:
		obj := common.NewImageObject(img)
		tex := common.NewTextureSingle(obj)
		s.videoEntity.Drawable.Close()
		s.videoEntity.Drawable = tex
	default:
		//do nothing
	}
}

func (s *VideoSystem) DataHandler(data []byte, wg *sync.WaitGroup) {
	// Create an image out of the data byte slice.
	img := image.NewNRGBA(
		image.Rectangle{
			image.Point{0, 0},
			image.Point{1280, 720},
		},
	)
	img.Pix = NRGBA(data)

	s.frameCh <- img

	wg.Done()
}

func NRGBA(rgbData []byte) []byte {
	numPixels := len(rgbData) / 3

	nrgbaData := make([]byte, numPixels*4)

	intNRGBAData := *(*[]uint32)(unsafe.Pointer(&nrgbaData))
	intNRGBAData = intNRGBAData[:len(nrgbaData)/4]

	for i, j := 0, 0; i < len(rgbData); i, j = i+3, j+1 {
		intRGB := (*(*uint32)(unsafe.Pointer(&rgbData[i]))) |
			(0b11111111 << 24)
		intNRGBAData[j] = intRGB
	}

	return nrgbaData
}
