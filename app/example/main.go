package main

import (
	"flag"

	"git.bug-br.org.br/bga/robomasters1/app"
)

var (
	ssId = flag.String("ssid", "testssid",
		"wifi network to connect to")
	password = flag.String("password", "testpassword", "wifi password")
	textMode = flag.Bool("textmode", false, "enable/disable text mode")
	appID    = flag.Uint64("appid", 0, "if provided, use this app ID "+
		"instead of creating a new one")
)

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
}
