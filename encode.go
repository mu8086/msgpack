package msgpack

import (
	"bytes"
	"encoding/binary"
)

func encode(buf *bytes.Buffer, data interface{}) error {
	switch v := data.(type) {
	case bool:
		return encodeBool(buf, v)

	// JSON numbers are unmarshaled as float64, so we need to handle int and uint in this case.
	case float64:
		if float64(uint64(v)) == v {
			return encodeUint(buf, uint64(v))
		} else if float64(int64(v)) == v {
			return encodeInt(buf, int64(v))
		}
		return encodeFloat(buf, v)

	case string:
		return encodeString(buf, v)

	default:
		return ErrUnsupportedType
	}
}

func encodeBool(buf *bytes.Buffer, value bool) error {
	// false (0xC2)
	if !value {
		return buf.WriteByte(0xC2)
	}
	// true (0xC3)
	return buf.WriteByte(0xC3)
}

func encodeFloat(_ *bytes.Buffer, _ float64) error {
	return nil
}

func encodeInt(buf *bytes.Buffer, value int64) error {
	switch {
	// positive fixint (0x00 ~ 0x7F)
	case value >= 0 && value <= 0x7F:
		buf.WriteByte(byte(value))

	// negative fixint (0xE0 ~ 0xFF)
	case value >= -32 && value <= -1:
		buf.WriteByte(0xE0 | byte(value+32))

	// int 8 (0xD0)
	case value >= -128 && value <= 127:
		buf.WriteByte(0xD0)
		buf.WriteByte(byte(value))

	// int 16 (0xD1)
	case value >= -32768 && value <= 32767:
		buf.WriteByte(0xD1)
		binary.Write(buf, binary.BigEndian, int16(value))

	// int 32 (0xD2)
	case value >= -2147483648 && value <= 2147483647:
		buf.WriteByte(0xD2)
		binary.Write(buf, binary.BigEndian, int32(value))

	// int 64 (0xD3)
	case value >= -9223372036854775808 && value <= 9223372036854775807:
		buf.WriteByte(0xD3)
		binary.Write(buf, binary.BigEndian, value)

	default:
		return ErrValueOutOfRange
	}

	return nil
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

func encodeUint(buf *bytes.Buffer, value uint64) error {
	switch {
	// positive fixint (0x00 ~ 0x7F)
	case value >= 0 && value <= 0x7F:
		buf.WriteByte(byte(value))

	// uint 8 (0xCC)
	case value <= 0xFF:
		buf.WriteByte(0xCC)
		buf.WriteByte(byte(value))

	// uint 16 (0xCD)
	case value <= 0xFFFF:
		buf.WriteByte(0xCD)
		binary.Write(buf, binary.BigEndian, uint16(value))

	// uint 32 (0xCE)
	case value <= 0xFFFFFFFF:
		buf.WriteByte(0xCE)
		binary.Write(buf, binary.BigEndian, uint32(value))

	// uint 64 (0xCF)
	case value <= 0xFFFFFFFFFFFFFFFF:
		buf.WriteByte(0xCF)
		binary.Write(buf, binary.BigEndian, value)

	default:
		return ErrValueOutOfRange
	}

	return nil
}
