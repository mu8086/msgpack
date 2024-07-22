package msgpack

import (
	"bytes"
	"encoding/json"

	"github.com/mu8086/msgpack/dto"
)

func JSONToMessagePack(jsonData []byte) ([]byte, error) {
	var data interface{}
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := encode(&buf, data); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func MessagePackToJSON(mp dto.MessagePack) ([]byte, error) {
	return mp.MarshalJSON()
}
