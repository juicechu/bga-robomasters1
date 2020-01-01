package main

import (
	"encoding/binary"
	"fmt"
	wrapper2 "git.bug-br.org.br/bga/robomasters1/app/internal/dji/unity/bridge/wrapper"
	"git.bug-br.org.br/bga/robomasters1/app/internal/dji/unity/bridge/wrapper/winebridge"
	"io"
	"os"
)

var (
	writePipe io.Writer

	readBuffer  []byte
	writeBuffer []byte

	wrapperInstance wrapper2.Wrapper
)

func main() {
	wineBridge, err := winebridge.New(2, os.Args[1:])
	if err != nil {
		panic(err)
	}

	readPipe := wineBridge.File(0)
	writePipe = wineBridge.File(1)

	// Initialize wrapper.
	wrapperInstance = wrapper2.Instance()

	lengthBuffer := make([]byte, 4)
	for {
		_, err := io.ReadFull(readPipe, lengthBuffer)
		if err == io.ErrUnexpectedEOF {
			panic(err)
		}
		if err == io.EOF {
			break
		}

		length := binary.LittleEndian.Uint32(lengthBuffer)

		err = processRead(readPipe, int(length-4))
		if err != nil {
			panic(err)
		}
	}
}

func processRead(readPipe io.Reader, length int) error {
	if length > len(readBuffer) {
		readBuffer = make([]byte, length)
	}

	sizedRequestBuffer := readBuffer[:length]
	_, err := io.ReadFull(readPipe, sizedRequestBuffer)
	if err != nil {
		return err
	}

	function := sizedRequestBuffer[0]
	switch function {
	case wrapper2.FuncCreateUnityBridge:
		runCreateUnityBridge(sizedRequestBuffer[1:])
	case wrapper2.FuncDestroyUnityBridge:
		runDestroyUnityBridge()
	case wrapper2.FuncUnityBridgeInitialize:
		runUnityBridgeInitialize()
	case wrapper2.FuncUnityBridgeUninitialize:
		runUnityBridgeUninitialize()
	case wrapper2.FuncUnitySetEventCallback:
		runUnitySetEventCallback(sizedRequestBuffer[1:])
	case wrapper2.FuncUnitySendEvent:
		runUnitySendEvent(sizedRequestBuffer[1:])
	case wrapper2.FuncUnitySendEventWithNumber:
		runUnitySendEventWithNumber(sizedRequestBuffer[1:])
	case wrapper2.FuncUnitySendEventWithString:
		runUnitySendEventWithString(sizedRequestBuffer[1:])
	}

	return nil
}

func runCreateUnityBridge(buffer []byte) {
	debuggable := false
	if buffer[0] != 0 {
		debuggable = true
	}

	name := string(buffer[1:])

	wrapperInstance.CreateUnityBridge(name, debuggable)
}

func runDestroyUnityBridge() {
	wrapperInstance.DestroyUnityBridge()
}

func runUnityBridgeInitialize() {
	if wrapperInstance.UnityBridgeInitialize() {
		fmt.Println("Unity Bridge initialized.")
	} else {
		fmt.Println("Unity Bridge failed to initialize.")
	}
}

func runUnityBridgeUninitialize() {
	wrapperInstance.UnityBridgeUninitialize()
}

func runUnitySetEventCallback(buffer []byte) {
	add := false
	if buffer[0] == 1 {
		add = true
	}
	eventCode := binary.LittleEndian.Uint64(buffer[1:])
	if add {
		wrapperInstance.UnitySetEventCallback(eventCode,
			eventCallback)
	} else {
		wrapperInstance.UnitySetEventCallback(eventCode, nil)
	}
}

func runUnitySendEvent(buffer []byte) {
	eventCode := binary.LittleEndian.Uint64(buffer)
	tag := binary.LittleEndian.Uint64(buffer[8:])
	info := buffer[16:]

	wrapperInstance.UnitySendEvent(eventCode, info, tag)
}

func runUnitySendEventWithNumber(buffer []byte) {
	eventCode := binary.LittleEndian.Uint64(buffer)
	tag := binary.LittleEndian.Uint64(buffer[8:])
	info := binary.LittleEndian.Uint64(buffer[16:])

	wrapperInstance.UnitySendEventWithNumber(eventCode, info, tag)
}

func runUnitySendEventWithString(buffer []byte) {
	eventCode := binary.LittleEndian.Uint64(buffer)
	tag := binary.LittleEndian.Uint64(buffer[8:])
	info := string(buffer[16:])

	wrapperInstance.UnitySendEventWithString(eventCode, info, tag)
}

func eventCallback(eventCode uint64, info []byte, tag uint64) {
	length := 4 + 8 + 8 + len(info)
	if len(writeBuffer) < length {
		writeBuffer = make([]byte, length)
	}

	sizedWriteBuffer := writeBuffer[:length]
	binary.LittleEndian.PutUint32(sizedWriteBuffer, uint32(length))
	binary.LittleEndian.PutUint64(sizedWriteBuffer[4:], eventCode)
	binary.LittleEndian.PutUint64(sizedWriteBuffer[12:], tag)
	copy(sizedWriteBuffer[20:], info)

	_, err := writePipe.Write(sizedWriteBuffer)
	if err != nil {
		panic(err)
	}
}
