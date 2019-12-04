package main

import (
	"flag"

	"git.bug-br.org.br/bga/robomasters1/app"
)

var (
	ssId     = flag.String("ssid", "testssid", "wifi network to connect to")
	password = flag.String("password", "testpassword", "wifi password")
	textMode = flag.Bool("textmode", false, "enable/disable text mode")
)

func main() {
	flag.Parse()

	a, err := app.New("US", *ssId, *password /*bssId=*/, "")
	if err != nil {
		panic(err)
	}

	err = a.Start(*textMode)
	if err != nil {
		panic(err)
	}
}
