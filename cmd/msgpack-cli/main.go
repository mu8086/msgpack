package main

import (
	"fmt"

	"github.com/mu8086/msgpack"
)

func main() {
	jsonData := []byte(`""`)

	mp, err := msgpack.JSONToMessagePack(jsonData)
	fmt.Printf("mp: %v (%v), err: %v\n", formatBytecode(mp), string(mp), err)
}

// TODO: remove
func formatBytecode(bytecode []byte) (s string) {
	for _, b := range bytecode {
		s += fmt.Sprintf("% #02x", b)
	}
	if len(s) > 0 { // remove leading space
		s = s[1:]
	}
	return s
}
