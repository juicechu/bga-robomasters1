package main

import (
	"flag"
	"image"
	"sync"
	"time"

	"git.bug-br.org.br/bga/robomasters1/app"
	"git.bug-br.org.br/bga/robomasters1/app/controller"
	"git.bug-br.org.br/bga/robomasters1/app/internal/rgb"
	"git.bug-br.org.br/bga/robomasters1/app/video"

	"fyne.io/fyne"
	fyneapp "fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
)

var (
	ssId = flag.String("ssid", "testssid",
		"wifi network to connect to")
	password = flag.String("password", "testpassword", "wifi password")
	textMode = flag.Bool("textmode", false, "enable/disable text mode")
	appID    = flag.Uint64("appid", 0, "if provided, use this app ID "+
		"instead of creating a new one")

	baseImg *rgb.Image
	img     *canvas.Image
)

func main() {
	flag.Parse()

	fyneApp := fyneapp.New()
	baseImg = rgb.NewImage(
		image.Rectangle{
			image.Point{0, 0},
			image.Point{1280, 720},
		},
	)
	img = canvas.NewImageFromImage(baseImg)
	img.FillMode = canvas.ImageFillOriginal

	w := fyneApp.NewWindow("Robomaster S1")
	w.Resize(fyne.Size{
		Width:  1280,
		Height: 720,
	})
	w.CenterOnScreen()
	w.SetContent(img)

	var a *app.App
	var err error
	if *appID != 0 {
		a, err = app.NewWithAppID("US", *ssId, *password,
			/*bssId=*/ "", *appID)
	} else {
		a, err = app.New("US", *ssId, *password /*bssId=*/, "")
	}
	if err != nil {
		panic(err)
	}

	err = a.Start(*textMode)
	if err != nil {
		panic(err)
	}

	v, err := video.New()
	if err != nil {
		panic(err)
	}

	index, err := v.AddDataHandler(videoHandler)
	if err != nil {
		panic(err)
	}
	defer v.RemoveDataHandler(index)

	v.StartVideo()

	c := controller.New(a.CommandController())

	// TODO(bga): HACK, fix me.
	time.Sleep(5 * time.Second)

	go func() {
		// Move the gimbal around.
		for i := 0; i < 100; i++ {
			c.Move(0.0, 0.0, 0.0, 1.0, false, true, 0)
		}
	}()

	w.ShowAndRun()
}

func videoHandler(data []byte, wg *sync.WaitGroup) {
	baseImg.Pix = data
	img.Refresh()

	wg.Done()
}
