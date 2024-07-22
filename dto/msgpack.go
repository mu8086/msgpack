package dto

import (
	"fmt"
)

type MessagePack struct {
	Bytecode []byte
}

// TODO:
// MessagePack to JSON
func (mp MessagePack) MarshalJSON() ([]byte, error) {
	size := len(mp.Bytecode)

	jsonData := make([]byte, size, size)
	copy(jsonData, mp.Bytecode)

	jsonData = append(jsonData, []byte("JSON")...)

	return jsonData, nil
}

func (mp MessagePack) String() (s string) {
	for _, b := range mp.Bytecode {
		s += fmt.Sprintf("% #02x", b)
	}
	if len(s) > 0 { // remove leading space
		s = s[1:]
	}
	return s
}

// TODO:
// JSON to MessagePack
func (mp *MessagePack) UnmarshalJSON(jsonData []byte) error {
	mp.Bytecode = append(mp.Bytecode, jsonData...)
	return nil
}
