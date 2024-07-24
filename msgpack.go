package msgpack

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func JSONToMessagePack(jsonData []byte) ([]byte, error) {
	var data interface{}
	if err := json.Unmarshal(jsonData, &data); err != nil {
		fmt.Printf("Unmarshal err: %v\n", err)
		return nil, err
	}

	var buf bytes.Buffer
	if err := encode(&buf, data); err != nil {
		fmt.Printf("encode err: %v\n", err)
		return nil, err
	}

	return buf.Bytes(), nil
}

func MessagePackToJSON(mp []byte) (string, error) {
	decoder := NewMessagePackDecoder(mp)

	data, err := decoder.Decode()
	if err != nil {
		return "", fmt.Errorf("MessagePack decode error: %v", err)
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("JSON marshal error: %v", err)
	}
	return string(jsonData), nil
}
