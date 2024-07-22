package msgpack

import (
	"github.com/mu8086/msgpack/dto"
)

func JSONToMessagePack(jsonData []byte) (mp dto.MessagePack) {
	mp.UnmarshalJSON(jsonData)
	return mp
}

func MessagePackToJSON(mp dto.MessagePack) ([]byte, error) {
	return mp.MarshalJSON()
}
