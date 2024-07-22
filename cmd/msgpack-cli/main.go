package main

import (
	"fmt"

	"github.com/mu8086/msgpack/pkg/msgpack"
)

func main() {
	jsonData := []byte("jsonData")

	mp := msgpack.JSONToMessagePack(jsonData)
	fmt.Printf("mp: \"%v\" (%v)\n", mp, string(mp.Bytecode))

	jsonData2, err := msgpack.MessagePackToJSON(mp)
	fmt.Printf("jsonData2: \"%v\" (%v), err: %v\n", jsonData2, string(jsonData2), err)
}
