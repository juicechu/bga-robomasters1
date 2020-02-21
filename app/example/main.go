package main

import (
	"flag"
	"image/color"

	"git.bug-br.org.br/bga/robomasters1/app"
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

var (
	ssId = flag.String("ssid", "testssid",
		"wifi network to connect to")
	password = flag.String("password", "testpassword", "wifi password")
	textMode = flag.Bool("textmode", false, "enable/disable text mode")
	appID    = flag.Uint64("appid", 0, "if provided, use this app ID "+
		"instead of creating a new one")
)

type DefaultScene struct{}

func (*DefaultScene) Preload() {
}

func (scene *DefaultScene) Setup(u engo.Updater) {
	w, _ := u.(*ecs.World)

	common.SetBackground(color.RGBA{R: 0, G: 255, B: 0, A: 255})

	w.AddSystem(&common.RenderSystem{})
	w.AddSystem(&VideoSystem{})
	w.AddSystem(&common.FPSSystem{
		Display: true,
	})
}

func (*DefaultScene) Type() string { return "RobomasterS1" }

func main() {
	flag.Parse()

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

	opts := engo.RunOptions{
		Title:  "Robomaster S1",
		Width:  1280,
		Height: 720,
	}
	engo.Run(opts, &DefaultScene{})
}
