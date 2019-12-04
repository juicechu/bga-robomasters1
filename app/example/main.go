package main

import (
	"flag"
	"fmt"

	"git.bug-br.org.br/bga/robomasters1/app"
)

var (
	ssId     = flag.String("ssid", "testssid", "wifi network to connect to")
	password = flag.String("password", "testpassword", "wifi password")
)

func main() {
	a, err := app.New("US", *ssId, *password /*bssId=*/, "")
	if err != nil {
		panic(err)
	}

	err = a.Start()
	if err != nil {
		panic(err)
	}
}
