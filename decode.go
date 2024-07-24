package msgpack

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type MessagePackDecoder struct {
	reader *bytes.Reader
}

func NewMessagePackDecoder(data []byte) *MessagePackDecoder {
	return &MessagePackDecoder{
		reader: bytes.NewReader(data),
	}
}

func (dec *MessagePackDecoder) Decode() (interface{}, error) {
	tag := "[MessagePackDecoder.Decode]"

	b, err := dec.reader.ReadByte()
	if err != nil {
		fmt.Printf("%v ReadByte failed, err: %v\n", tag, err)
		return nil, err
	}

	switch {
	// positive fixint
	case b >= 0x00 && b <= 0x7F:
		return uint8(b), nil

	// fixmap
	case b >= 0x80 && b <= 0x8F:
		length := int(b & 0x0F)
		return dec.readMap(length)

	// fixarray
	case b >= 0x90 && b <= 0x9F:
		length := int(b & 0x0F)
		return dec.readArray(length)

	// fixstr
	case b >= 0xA0 && b <= 0xBF:
		length := int(b & 0x1F)
		return dec.readString(length)

	// nil
	case b == 0xC0:
		return nil, nil

	// never used
	case b == 0xC1:
		fmt.Printf("%v 0x%02X not used in MessagePack\n", tag, b)
		return nil, ErrUnsupportedType

	// false
	case b == 0xC2:
		return false, nil

	// true
	case b == 0xC3:
		return true, nil

	// bin 8
	case b == 0xC4:
		return dec.readBinWithLengthInBits(8)

	// bin 16
	case b == 0xC5:
		return dec.readBinWithLengthInBits(16)

	// bin 32
	case b == 0xC6:
		return dec.readBinWithLengthInBits(32)

	// ext 8
	case b == 0xC7:

	// ext 16
	case b == 0xC8:

	// ext 32
	case b == 0xC9:

	// float 32
	case b == 0xCA:
		return dec.readFloat32()

	// float 64
	case b == 0xCB:
		return dec.readFloat64()

	// uint8
	case b == 0xCC:
		return dec.readUint8()

	// uint 16
	case b == 0xCD:
		return dec.readUint16()

	// uint 32
	case b == 0xCE:
		return dec.readUint32()

	// uint 64
	case b == 0xCF:
		return dec.readUint64()

	// int 8
	case b == 0xD0:
		return dec.readInt8()

	// int 16
	case b == 0xD1:
		return dec.readInt16()

	// int 32
	case b == 0xD2:
		return dec.readInt32()

	// int 64
	case b == 0xD3:
		return dec.readInt64()

	// fixext 1
	case b == 0xD4:

	// fixext 2
	case b == 0xD5:

	// fixext 4
	case b == 0xD6:

	// fixext 8
	case b == 0xD7:

	// fixext 16
	case b == 0xD8:

	// str 8
	case b == 0xD9:
		return dec.readStrWithLengthInBits(8)

	// str 16
	case b == 0xDA:
		return dec.readStrWithLengthInBits(16)

	// str 32
	case b == 0xDB:
		return dec.readStrWithLengthInBits(32)

	// array 16
	case b == 0xDC:
		return dec.readArrayWithLengthInBits(16)

	// array 32
	case b == 0xDD:
		return dec.readArrayWithLengthInBits(32)

	// map 16
	case b == 0xDE:
		return dec.readMapWithLengthInBits(16)

	// map 32
	case b == 0xDF:
		return dec.readMapWithLengthInBits(32)

	// negative fixint
	case b >= 0xE0 && b <= 0xFF:
		return int8(b), nil

	default:
		fmt.Printf("%v 0x%02X not defined in MessagePack\n", tag, b)
		return "", ErrUnsupportedType
	}

	return "", nil
}

func (dec *MessagePackDecoder) readArray(length int) ([]interface{}, error) {
	tag := "[MessagePackDecoder.readArray]"

	data := make([]interface{}, length, length)

	for i := 0; i < int(length); i++ {
		element, err := dec.Decode()
		if err != nil {
			fmt.Printf("%v Decode failed, err: %v\n", tag, err)
			return nil, err
		}

		data[i] = element
	}

	return data, nil
}

func (dec *MessagePackDecoder) readArrayWithLengthInBits(lengthInBits int) ([]interface{}, error) {
	tag := "[MessagePackDecoder.readArrayWithLengthInBits]"

	length, err := dec.readLength(lengthInBits)
	if err != nil {
		fmt.Printf("%v readLength failed, err: %v\n", tag, err)
		return nil, err
	}

	return dec.readArray(int(length))
}

func (dec *MessagePackDecoder) readBinWithLengthInBits(lengthInBits int) ([]byte, error) {
	tag := "[MessagePackDecoder.readBinWithLengthInBits]"

	length, err := dec.readLength(lengthInBits)
	if err != nil {
		fmt.Printf("%v readLength failed, err: %v\n", tag, err)
		return nil, err
	}

	binData := make([]byte, length, length)
	err = binary.Read(dec.reader, binary.BigEndian, &binData)
	if err != nil {
		fmt.Printf("%v Read failed, err: %v\n", tag, err)
		return nil, err
	}

	return binData, nil
}

func (dec *MessagePackDecoder) readFloat32() (data float32, err error) {
	err = binary.Read(dec.reader, binary.BigEndian, &data)
	return data, err
}

func (dec *MessagePackDecoder) readFloat64() (data float64, err error) {
	err = binary.Read(dec.reader, binary.BigEndian, &data)
	return data, err
}

func (dec *MessagePackDecoder) readInt8() (data int8, err error) {
	err = binary.Read(dec.reader, binary.BigEndian, &data)
	return data, err
}

func (dec *MessagePackDecoder) readInt16() (data int16, err error) {
	err = binary.Read(dec.reader, binary.BigEndian, &data)
	return data, err
}

func (dec *MessagePackDecoder) readInt32() (data int32, err error) {
	err = binary.Read(dec.reader, binary.BigEndian, &data)
	return data, err
}

func (dec *MessagePackDecoder) readInt64() (data int64, err error) {
	err = binary.Read(dec.reader, binary.BigEndian, &data)
	return data, err
}

func (dec *MessagePackDecoder) readLength(bits int) (int64, error) {
	tag := "[MessagePackDecoder.readLength]"

	// not a multiple of 8
	if (bits & 0x7) != 0 {
		fmt.Printf("%v bits(%v) not a multiple of 8\n", tag, bits)
		return 0, ErrLengthInvalid
	}

	length := int64(0)
	byteSize := bits >> 3

	for i := 1; i <= byteSize; i++ {
		length <<= 4

		b, err := dec.reader.ReadByte()
		if err != nil {
			fmt.Printf("%v ReadByte failed, err: %v\n", tag, err)
			return 0, ErrReadByte
		}

		length |= int64(b)
	}

	return length, nil
}

func (dec *MessagePackDecoder) readMap(length int) (map[string]interface{}, error) {
	tag := "[MessagePackDecoder.readMap]"

	m := make(map[string]interface{}, length)

	for i := 0; i < length; i++ {
		key, err := dec.Decode()
		if err != nil {
			fmt.Printf("%v Decode key failed, err: %v\n", tag, err)
			return nil, err
		}

		keyStr, ok := key.(string)
		if !ok {
			fmt.Printf("%v key not a string\n", tag)
			return nil, ErrUnsupportedType
		}

		value, err := dec.Decode()
		if err != nil {
			fmt.Printf("%v Decode value failed, err: %v\n", tag, err)
			return nil, err
		}

		m[keyStr] = value
	}

	return m, nil
}

func (dec *MessagePackDecoder) readMapWithLengthInBits(lengthInBits int) (map[string]interface{}, error) {
	tag := "[MessagePackDecoder.readMapWithLengthInBits]"

	length, err := dec.readLength(lengthInBits)
	if err != nil {
		fmt.Printf("%v readLength failed, err: %v\n", tag, err)
		return nil, err
	}

	return dec.readMap(int(length))
}

func (dec *MessagePackDecoder) readStrWithLengthInBits(lengthInBits int) (string, error) {
	tag := "[MessagePackDecoder.readStrWithLengthInBits]"

	length, err := dec.readLength(lengthInBits)
	if err != nil {
		fmt.Printf("%v readLength failed, err: %v\n", tag, err)
		return "", err
	}

	return dec.readString(int(length))
}

func (dec *MessagePackDecoder) readString(length int) (string, error) {
	tag := "[MessagePackDecoder.readString]"

	buf := make([]byte, length)
	if _, err := dec.reader.Read(buf); err != nil {
		fmt.Printf("%v Read failed, err: %v\n", tag, err)
		return "", err
	}
	return string(buf), nil
}

func (dec *MessagePackDecoder) readUint8() (data uint8, err error) {
	tag := "[MessagePackDecoder.readUint8]"

	err = binary.Read(dec.reader, binary.BigEndian, &data)
	if err != nil {
		fmt.Printf("%v Read failed, err: %v\n", tag, err)
		return 0, err
	}
	return data, nil
}

func (dec *MessagePackDecoder) readUint16() (data uint16, err error) {
	tag := "[MessagePackDecoder.readUint16]"

	err = binary.Read(dec.reader, binary.BigEndian, &data)
	if err != nil {
		fmt.Printf("%v Read failed, err: %v\n", tag, err)
		return 0, err
	}
	return data, err
}

func (dec *MessagePackDecoder) readUint32() (data uint32, err error) {
	tag := "[MessagePackDecoder.readUint32]"

	err = binary.Read(dec.reader, binary.BigEndian, &data)
	if err != nil {
		fmt.Printf("%v Read failed, err: %v\n", tag, err)
		return 0, err
	}
	return data, err
}

func (dec *MessagePackDecoder) readUint64() (data uint64, err error) {
	tag := "[MessagePackDecoder.readUint64]"

	err = binary.Read(dec.reader, binary.BigEndian, &data)
	if err != nil {
		fmt.Printf("%v Read failed, err: %v\n", tag, err)
		return 0, err
	}
	return data, err
}
