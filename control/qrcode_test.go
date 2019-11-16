package control

import (
	"testing"
)

func TestParseQRCodeMessage(t *testing.T) {
	q, err := ParseQRCodeMessage("jr+7eFsfyQBLc9xx/3nWfeBD81aXMorTfcNq8UPuVPtd+WM=")
	if err != nil {
		t.Fatalf("expected nil error, got %q", err)
	}

	expected := "App Id : 5130944286501155530, Country Code : \"US\", SSID : \"discworld\", Password : \"zwergschnauzer\", BSSID : <nil>"
	if q.String() != expected {
		t.Fatalf("expected %q, got %q", expected, q.String())
	}
}
