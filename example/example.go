package main

import (
	"encoding/binary"
	"fmt"
	"time"

	"github.com/brunoga/unitybridge"
	"github.com/brunoga/unitybridge/event"
)

func callbackHandler(e *event.Event, data []byte, tag uint64) {
	dataType := (tag >> 56) & 0xff

	dataTypeStr := "unknown"
	switch dataType {
	case 0:
		dataTypeStr = "string"
	case 1:
		dataTypeStr = "uint64"
	}

	tag = tag & 0x0000ffffffffffff

	fmt.Printf("Callback handler called for event with type %s, sub-type %d, data type %s and tag %d\n",
		e.Type(), e.SubType(), dataTypeStr, tag)

	if dataType == 0 {
		fmt.Printf("Data: %s\n", string(data))
	} else {
		fmt.Printf("Data: %d\n", binary.NativeEndian.Uint64(data))
	}
}

func main() {
	ub := unitybridge.Get()

	ub.Create("Robomaster", true, "./log")
	defer ub.Destroy()

	if !ub.Initialize() {
		panic("Could not initialize UnityBridge.")
	}
	defer ub.Uninitialize()

	for _, typ := range event.AllTypes() {
		ub.SetEventCallback(typ, callbackHandler)
	}

	time.Sleep(5 * time.Second)

	for _, typ := range event.AllTypes() {
		ub.SetEventCallback(typ, nil)
	}
}
