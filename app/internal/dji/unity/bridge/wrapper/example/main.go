package main

import (
	"fmt"
	"runtime"
	"time"

	"git.bug-br.org.br/bga/robomasters1/app/internal/dji/unity/bridge/wrapper"
)

func main() {
	w, err := wrapper.New("./unitybridge.dll")
	if err != nil {
		panic(err)
	}

	w.CreateUnityBridge("Robomaster", true)

	events := []uint64{
		0, 1, 2, 3, 4, 5, 6, 7, 8, 100, 101, 200, 300, 301, 302, 303,
		304, 305, 306, 500,
	}
	for _, event := range events {
		w.UnitySetEventCallback(event<<32, callback)
	}

	if !w.UnityBridgeInitialize() {
		panic("Failed initializing unity bridge")
	}

	time.Sleep(10 * time.Second)

	w.UnityBridgeUninitialize()

	for _, event := range events {
		w.UnitySetEventCallback(event<<32, nil)
	}

	w.DestroyUnityBridge()

	w = nil
	runtime.GC()

	time.Sleep(5 * time.Second)
}

func callback(eventCode uint64, data []byte, tag uint64) {
	fmt.Printf("callback(eventCode=%v, data=%v, tag=%v)\n", eventCode, data,
		tag)
}
