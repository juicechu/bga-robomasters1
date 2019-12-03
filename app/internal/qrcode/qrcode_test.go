package qrcode

import (
	"testing"
)

func TestParseQRCodeMessage(t *testing.T) {
	q, err := ParseQRCodeMessage("A767eFsfyQBLc9xx6GPMeudN8kmEJ4/S")
	if err != nil {
		t.Fatalf("expected nil error, got %q", err)
	}

	expected := "App Id : 5130944286501155530, Country Code : \"US\", " +
		"SSID : \"ssid\", Password : \"password\", BSSID : <nil>"
	if q.String() != expected {
		t.Fatalf("expected %q, got %q", expected, q.String())
	}
}
