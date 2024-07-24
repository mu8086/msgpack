package msgpack

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"math"
)

func encode(buf *bytes.Buffer, data interface{}) error {
	tag := "[encode]"

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

	case int:
		return encodeInt(buf, int64(v))

	case int64:
		return encodeInt(buf, v)

	case []interface{}:
		return encodeArray(buf, v)

	case map[string]interface{}:
		return encodeMap(buf, v)

	case nil:
		return encodeNil(buf, v)

	case string:
		return encodeString(buf, v)

	case uint:
		return encodeUint(buf, uint64(v))

	case uint64:
		return encodeUint(buf, v)

	default:
		fmt.Printf("%v Unsupported Type: %T\n", tag, v)
		return ErrUnsupportedType
	}
}

func encodeArray(buf *bytes.Buffer, value []interface{}) error {
	tag := "[encodeArray]"

	length := len(value)

	switch {
	// fixarray (0x90 ~ 0x9F)
	case length <= 0xF:
		buf.WriteByte(0x90 | byte(length))

	// array 16 (0xDC)
	case length <= 0xFFFF:
		buf.WriteByte(0xDC)
		binary.Write(buf, binary.BigEndian, int16(length))

	// array 32 (0xDD)
	case length <= 0xFFFFFFFF:
		buf.WriteByte(0xDD)
		binary.Write(buf, binary.BigEndian, int32(length))

	default:
		fmt.Printf("%v array size (%v) too large\n", tag, length)
		return ErrArrayTooLong
	}

	for _, element := range value {
		if err := encode(buf, element); err != nil {
			fmt.Printf("%v encode failed, err: %v\n", tag, err)
			return err
		}
	}

	return nil
}

func encodeBinary(buf *bytes.Buffer, value interface{}) error {
	tag := "[encodeBinary]"

	base64Str, ok := value.(string)
	if !ok {
		return ErrBinaryDataInvalid
	}

	binData, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		fmt.Printf("%v DecodeString failed, err: %v\n", tag, err)
		return err
	}

	length := len(binData)

	switch {
	// bin 8 (0xC4)
	case length <= 0xFF: // 2^8 - 1
		err = buf.WriteByte(0xC4)
		if err != nil {
			fmt.Printf("%v WriteByte failed, err: %v\n", tag, err)
			return err
		}
		err = binary.Write(buf, binary.BigEndian, int8(length))
		if err != nil {
			fmt.Printf("%v Write failed, err: %v\n", tag, err)
			return err
		}

	// bin 16 (0xC5)
	case length <= 0xFFFF: // 2^16 - 1
		err = buf.WriteByte(0xC5)
		if err != nil {
			fmt.Printf("%v WriteByte failed, err: %v\n", tag, err)
			return err
		}
		err = binary.Write(buf, binary.BigEndian, int16(length))
		if err != nil {
			fmt.Printf("%v Write failed, err: %v\n", tag, err)
			return err
		}

	// bin 32 (0xC6)
	case length <= 0xFFFFFFFF: // 2^32 - 1
		err = buf.WriteByte(0xC6)
		if err != nil {
			fmt.Printf("%v WriteByte failed, err: %v\n", tag, err)
			return err
		}
		err = binary.Write(buf, binary.BigEndian, int32(length))
		if err != nil {
			fmt.Printf("%v Write failed, err: %v\n", tag, err)
			return err
		}

	default:
		fmt.Printf("%v bin size (%v) too large\n", tag, length)
		return ErrBinaryTooLong
	}

	_, err = buf.Write(binData)
	if err != nil {
		fmt.Printf("%v Write failed, err: %v\n", tag, err)
		return err
	}
	return nil
}

func encodeBool(buf *bytes.Buffer, value bool) error {
	// false (0xC2)
	if !value {
		return buf.WriteByte(0xC2)
	}
	// true (0xC3)
	return buf.WriteByte(0xC3)
}

func encodeFloat(buf *bytes.Buffer, value float64) error {
	tag := "[encodeFloat]"

	if float64(float32(value)) == value {
		// float 32 (0xCA)
		err := buf.WriteByte(0xCA)
		if err != nil {
			fmt.Printf("%v WriteByte failed, err: %v\n", tag, err)
			return err
		}

		bits := math.Float32bits(float32(value))

		err = binary.Write(buf, binary.BigEndian, bits)
		if err != nil {
			fmt.Printf("%v Write failed, err: %v\n", tag, err)
			return err
		}
	} else {
		// float 64 (0xCB)
		err := buf.WriteByte(0xCB)
		if err != nil {
			fmt.Printf("%v WriteByte failed, err: %v\n", tag, err)
			return err
		}

		bits := math.Float64bits(value)

		err = binary.Write(buf, binary.BigEndian, bits)
		if err != nil {
			fmt.Printf("%v Write failed, err: %v\n", tag, err)
			return err
		}
	}
	return nil
}

func encodeInt(buf *bytes.Buffer, value int64) (err error) {
	tag := "[encodeInt]"

	switch {
	// positive fixint (0x00 ~ 0x7F)
	case value >= 0 && value <= 0x7F:
		err = buf.WriteByte(byte(value))
		if err != nil {
			fmt.Printf("%v WriteByte failed, err: %v\n", tag, err)
			return err
		}

	// negative fixint (0xE0 ~ 0xFF)
	case value >= -32 && value <= -1:
		err = buf.WriteByte(0xE0 | byte(value+32))
		if err != nil {
			fmt.Printf("%v WriteByte failed, err: %v\n", tag, err)
			return err
		}

	// int 8 (0xD0)
	case value >= -128 && value <= 127:
		err = buf.WriteByte(0xD0)
		if err != nil {
			fmt.Printf("%v WriteByte failed, err: %v\n", tag, err)
			return err
		}

		err = buf.WriteByte(byte(value))
		if err != nil {
			fmt.Printf("%v WriteByte failed, err: %v\n", tag, err)
			return err
		}

	// int 16 (0xD1)
	case value >= -32768 && value <= 32767:
		err = buf.WriteByte(0xD1)
		if err != nil {
			fmt.Printf("%v WriteByte failed, err: %v\n", tag, err)
			return err
		}

		err = binary.Write(buf, binary.BigEndian, int16(value))
		if err != nil {
			fmt.Printf("%v Write failed, err: %v\n", tag, err)
			return err
		}

	// int 32 (0xD2)
	case value >= -2147483648 && value <= 2147483647:
		err = buf.WriteByte(0xD2)
		if err != nil {
			fmt.Printf("%v WriteByte failed, err: %v\n", tag, err)
			return err
		}

		err = binary.Write(buf, binary.BigEndian, int32(value))
		if err != nil {
			fmt.Printf("%v Write failed, err: %v\n", tag, err)
			return err
		}

	// int 64 (0xD3)
	case value >= -9223372036854775808 && value <= 9223372036854775807:
		err = buf.WriteByte(0xD3)
		if err != nil {
			fmt.Printf("%v WriteByte failed, err: %v\n", tag, err)
			return err
		}

		err = binary.Write(buf, binary.BigEndian, value)
		if err != nil {
			fmt.Printf("%v Write failed, err: %v\n", tag, err)
			return err
		}

	default:
		fmt.Printf("%v not supported value(%v)\n", tag, value)
		return ErrValueOutOfRange
	}

	return nil
}

func encodeMap(buf *bytes.Buffer, value map[string]interface{}) (err error) {
	tag := "[encodeMap]"

	length := len(value)

	switch {
	//fixmap (0x80 ~ 0x8F)
	case length <= 0xF:
		err = buf.WriteByte(0x80 | byte(length))
		if err != nil {
			fmt.Printf("%v WriteByte failed, err: %v\n", tag, err)
			return err
		}

	// map 16 (0xDE)
	case length <= 0xFFFF:
		err = buf.WriteByte(0xDE)
		if err != nil {
			fmt.Printf("%v WriteByte failed, err: %v\n", tag, err)
			return err
		}

		err = binary.Write(buf, binary.BigEndian, int16(length))
		if err != nil {
			fmt.Printf("%v Write failed, err: %v\n", tag, err)
			return err
		}

	// map 32 (0xDF)
	case length <= 0xFFFFFFFF:
		err = buf.WriteByte(0xDF)
		if err != nil {
			fmt.Printf("%v WriteByte failed, err: %v\n", tag, err)
			return err
		}

		err = binary.Write(buf, binary.BigEndian, int32(length))
		if err != nil {
			fmt.Printf("%v Write failed, err: %v\n", tag, err)
			return err
		}

	default:
		fmt.Printf("%v map size(%v) too large\n", tag, length)
		return ErrValueOutOfRange
	}

	for key, val := range value {
		if err := encodeString(buf, key); err != nil {
			fmt.Printf("%v key encodeString failed, err: %v\n", tag, err)
			return err
		}

		if key == binaryKeyword {
			if err := encodeBinary(buf, val); err != nil {
				fmt.Printf("%v encodeBinary failed, err: %v\n", tag, err)
				return err
			}
		} else {
			if err := encode(buf, val); err != nil {
				fmt.Printf("%v value encode failed, err: %v\n", tag, err)
				return err
			}
		}
	}
	return nil
}

func encodeNil(buf *bytes.Buffer, _ interface{}) error {
	return buf.WriteByte(0xC0)
}

func encodeString(buf *bytes.Buffer, value string) (err error) {
	tag := "[encodeString]"

	length := len(value)

	switch {
	// fixstr (0xA0 - 0xBF)
	case length <= 0x1F: // 31
		err = buf.WriteByte(0xA0 | byte(length))
		if err != nil {
			fmt.Printf("%v WriteByte failed, err: %v\n", tag, err)
			return err
		}

	// str 8 (0xD9)
	case length <= 0xFF: // 2^8 - 1
		err = buf.WriteByte(0xD9)
		if err != nil {
			fmt.Printf("%v WriteByte failed, err: %v\n", tag, err)
			return err
		}

		err = buf.WriteByte(byte(length))
		if err != nil {
			fmt.Printf("%v WriteByte failed, err: %v\n", tag, err)
			return err
		}

	// str 16 (0xDA)
	case length <= 0xFFFF: // 2^16 - 1
		err = buf.WriteByte(0xDA)
		if err != nil {
			fmt.Printf("%v WriteByte failed, err: %v\n", tag, err)
			return err
		}

		_, err = buf.Write([]byte{byte(length >> 8), byte(length)})
		if err != nil {
			fmt.Printf("%v Write failed, err: %v\n", tag, err)
			return err
		}

	// str 32 (0xDB)
	case length <= 0xFFFFFFFF: // 2^32 - 1
		err = buf.WriteByte(0xDB)
		if err != nil {
			fmt.Printf("%v WriteByte failed, err: %v\n", tag, err)
			return err
		}

		_, err = buf.Write([]byte{byte(length >> 24), byte(length >> 16), byte(length >> 8), byte(length)})
		if err != nil {
			fmt.Printf("%v Write failed, err: %v\n", tag, err)
			return err
		}

	default:
		fmt.Printf("%v string size(%v) too large\n", tag, length)
		return ErrStringTooLong
	}

	// write string content to buffer
	_, err = buf.WriteString(value)
	if err != nil {
		fmt.Printf("%v WriteString failed, err: %v\n", tag, err)
		return err
	}

	return nil
}

func encodeUint(buf *bytes.Buffer, value uint64) (err error) {
	tag := "[encodeUint]"

	switch {
	// positive fixint (0x00 ~ 0x7F)
	case value >= 0 && value <= 0x7F:
		err = buf.WriteByte(byte(value))
		if err != nil {
			fmt.Printf("%v WriteByte failed, err: %v\n", tag, err)
			return err
		}

	// uint 8 (0xCC)
	case value <= 0xFF:
		err = buf.WriteByte(0xCC)
		if err != nil {
			fmt.Printf("%v WriteByte failed, err: %v\n", tag, err)
			return err
		}

		err = buf.WriteByte(byte(value))
		if err != nil {
			fmt.Printf("%v WriteByte failed, err: %v\n", tag, err)
			return err
		}

	// uint 16 (0xCD)
	case value <= 0xFFFF:
		err = buf.WriteByte(0xCD)
		if err != nil {
			fmt.Printf("%v WriteByte failed, err: %v\n", tag, err)
			return err
		}

		err = binary.Write(buf, binary.BigEndian, uint16(value))
		if err != nil {
			fmt.Printf("%v Write failed, err: %v\n", tag, err)
			return err
		}

	// uint 32 (0xCE)
	case value <= 0xFFFFFFFF:
		err = buf.WriteByte(0xCE)
		if err != nil {
			fmt.Printf("%v WriteByte failed, err: %v\n", tag, err)
			return err
		}

		err = binary.Write(buf, binary.BigEndian, uint32(value))
		if err != nil {
			fmt.Printf("%v Write failed, err: %v\n", tag, err)
			return err
		}

	// uint 64 (0xCF)
	case value <= 0xFFFFFFFFFFFFFFFF:
		err = buf.WriteByte(0xCF)
		if err != nil {
			fmt.Printf("%v WriteByte failed, err: %v\n", tag, err)
			return err
		}

		err = binary.Write(buf, binary.BigEndian, value)
		if err != nil {
			fmt.Printf("%v Write failed, err: %v\n", tag, err)
			return err
		}

	default:
		fmt.Printf("%v not supported value(%v)\n", tag, value)
		return ErrValueOutOfRange
	}

	return nil
}
