package msgpack

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/mu8086/msgpack/dto"
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

func MessagePackToJSON(mp dto.MessagePack) ([]byte, error) {
	return mp.MarshalJSON()
}
