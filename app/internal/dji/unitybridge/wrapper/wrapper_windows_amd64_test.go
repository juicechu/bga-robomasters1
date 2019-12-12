package bridge

import (
	"testing"
)

func TestWrapperMinimalInitAndTearDown(t *testing.T) {
	w := Instance()
	w.CreateUnityBridge("Robomaster", true)
	defer w.DestroyUnityBridge()

	ok := w.UnityBridgeInitialize()
	if !ok {
		t.Fatal("cannot initialize unity bridge")
	}
	defer w.UnityBridgeUninitialize()
}
