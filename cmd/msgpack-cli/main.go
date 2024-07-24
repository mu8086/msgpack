package main

import (
	"fmt"

	"github.com/mu8086/msgpack"
)

func main() {
	jsonData := []byte(
		`{
			"0": true,
			"1": true,
			"2": true,
			"3": true,
			"4": true,
			"5": true,
			"6": true,
			"7": true,
			"8": true,
			"9": true,
			"A": true,
			"B": true,
			"C": true,
			"D": true,
			"E": true,
			"F": true
		}`)

	mp, err := msgpack.JSONToMessagePack(jsonData)
	fmt.Printf("mp: %v, err: %v\n", formatBytecode(mp), err)
}

// TODO: remove
func formatBytecode(bytecode []byte) (s string) {
	for _, b := range bytecode {
		s += fmt.Sprintf("% 03X", b)
	}
	if len(s) > 0 { // remove leading space
		s = s[1:]
	}
	return s
}
