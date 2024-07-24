package msgpack

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func JSONToMessagePack(jsonData []byte) ([]byte, error) {
	tag := "[JSONToMessagePack]"

	var data interface{}
	if err := json.Unmarshal(jsonData, &data); err != nil {
		fmt.Printf("%v Unmarshal failed, err: %v\n", tag, err)
		return nil, err
	}

	var buf bytes.Buffer
	if err := encode(&buf, data); err != nil {
		fmt.Printf("%v encode failed, err: %v\n", tag, err)
		return nil, err
	}

	return buf.Bytes(), nil
}

func MessagePackToJSON(mp []byte) (string, error) {
	tag := "[MessagePackToJSON]"

	decoder := NewMessagePackDecoder(mp)

	data, err := decoder.Decode()
	if err != nil {
		fmt.Printf("%v Decode failed, err: %v", tag, err)
		return "", err
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("%v Marshal failed, err: %v", tag, err)
		return "", err
	}
	return string(jsonData), nil
}
