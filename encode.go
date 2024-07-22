package msgpack

import (
	"bytes"
)

func encode(buf *bytes.Buffer, data interface{}) error {
	switch v := data.(type) {
	case bool:
		return encodeBool(buf, v)
	case string:
		return encodeString(buf, v)
	default:
		return ErrUnsupportedType
	}
}

func encodeBool(buf *bytes.Buffer, v bool) error {
	if !v {
		return buf.WriteByte(0xC2)
	}
	return buf.WriteByte(0xC3)
}

func encodeString(buf *bytes.Buffer, value string) error {
	length := len(value)

	switch {
	// fixstr (0xA0 - 0xBF)
	case length <= 0x1F: // 31
		buf.WriteByte(0xA0 | byte(length))

	// str 8 (0xD9)
	case length <= 0xFF: // 2^8 - 1
		buf.WriteByte(0xD9)
		buf.WriteByte(byte(length))

	// str 16 (0xDA)
	case length <= 0xFFFF: // 2^16 - 1
		buf.WriteByte(0xDA)
		buf.Write([]byte{byte(length >> 8), byte(length)})

	// str 32 (0xDB)
	case length <= 0xFFFFFFFF: // 2^32 - 1
		buf.WriteByte(0xDB)
		buf.Write([]byte{byte(length >> 24), byte(length >> 16), byte(length >> 8), byte(length)})

	default:
		return ErrStringTooLong
	}

	// write string content to buffer
	buf.WriteString(value)
	return nil
}
