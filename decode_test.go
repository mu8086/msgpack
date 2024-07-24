package msgpack

import (
	"bytes"
	"testing"
)

func TestPositiveFixInt(t *testing.T) {
	tests := []struct {
		input    []byte
		expected interface{}
	}{
		{input: []byte{0x00}, expected: uint8(0)},
		{input: []byte{0x7F}, expected: uint8(127)},
	}

	for _, test := range tests {
		decoder := NewMessagePackDecoder(test.input)
		result, err := decoder.Decode()
		if err != nil {
			t.Errorf("Decode() error = %v", err)
		}
		if !isEqual(result, test.expected) {
			t.Errorf("Decode() = %v, want %v", result, test.expected)
		}
	}
}

func TestFixMap(t *testing.T) {
	input := []byte{0x81, 0xA3, 'k', 'e', 'y', 0x01}
	expected := map[string]interface{}{"key": uint8(1)}

	decoder := NewMessagePackDecoder(input)
	result, err := decoder.Decode()
	if err != nil {
		t.Errorf("Decode() error = %v", err)
	}
	if !isEqual(result, expected) {
		t.Errorf("Decode() = %v, want %v", result, expected)
	}
}

func TestFixArray(t *testing.T) {
	input := []byte{0x92, 0x01, 0x02}
	expected := []interface{}{uint8(1), uint8(2)}

	decoder := NewMessagePackDecoder(input)
	result, err := decoder.Decode()
	if err != nil {
		t.Errorf("Decode() error = %v", err)
	}
	if !isEqual(result, expected) {
		t.Errorf("Decode() = %v, want %v", result, expected)
	}
}

func TestFixStr(t *testing.T) {
	input := []byte{0xA3, 'f', 'o', 'o'}
	expected := "foo"

	decoder := NewMessagePackDecoder(input)
	result, err := decoder.Decode()
	if err != nil {
		t.Errorf("Decode() error = %v", err)
	}
	if !isEqual(result, expected) {
		t.Errorf("Decode() = %v, want %v", result, expected)
	}
}

func TestNil(t *testing.T) {
	input := []byte{0xC0}
	//expected := nil

	decoder := NewMessagePackDecoder(input)
	result, err := decoder.Decode()
	if err != nil {
		t.Errorf("Decode() error = %v", err)
	}
	if !isEqual(result, nil) {
		t.Errorf("Decode() = %v, want %v", result, nil)
	}
}

func TestFalse(t *testing.T) {
	input := []byte{0xC2}
	expected := false

	decoder := NewMessagePackDecoder(input)
	result, err := decoder.Decode()
	if err != nil {
		t.Errorf("Decode() error = %v", err)
	}
	if !isEqual(result, expected) {
		t.Errorf("Decode() = %v, want %v", result, expected)
	}
}

func TestTrue(t *testing.T) {
	input := []byte{0xC3}
	expected := true

	decoder := NewMessagePackDecoder(input)
	result, err := decoder.Decode()
	if err != nil {
		t.Errorf("Decode() error = %v", err)
	}
	if !isEqual(result, expected) {
		t.Errorf("Decode() = %v, want %v", result, expected)
	}
}

func TestBin8(t *testing.T) {
	input := []byte{0xC4, 0x04, 0x01, 0x02, 0x03, 0x04}
	expected := []byte{0x01, 0x02, 0x03, 0x04}

	decoder := NewMessagePackDecoder(input)
	result, err := decoder.Decode()
	if err != nil {
		t.Errorf("Decode() error = %v", err)
	}
	if !isEqual(result, expected) {
		t.Errorf("Decode() = %v, want %v", result, expected)
	}
}

func TestFloat32(t *testing.T) {
	input := []byte{0xCA, 0x3F, 0x80, 0x00, 0x00}
	expected := float32(1.0)

	decoder := NewMessagePackDecoder(input)
	result, err := decoder.Decode()
	if err != nil {
		t.Errorf("Decode() error = %v", err)
	}
	if !isEqual(result, expected) {
		t.Errorf("Decode() = %v, want %v", result, expected)
	}
}

func TestFloat64(t *testing.T) {
	input := []byte{0xCB, 0x3F, 0xF0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	expected := float64(1.0)

	decoder := NewMessagePackDecoder(input)
	result, err := decoder.Decode()
	if err != nil {
		t.Errorf("Decode() error = %v", err)
	}
	if !isEqual(result, expected) {
		t.Errorf("Decode() = %v, want %v", result, expected)
	}
}

func TestUint8(t *testing.T) {
	input := []byte{0xCC, 0xFF}
	expected := uint8(255)

	decoder := NewMessagePackDecoder(input)
	result, err := decoder.Decode()
	if err != nil {
		t.Errorf("Decode() error = %v", err)
	}
	if !isEqual(result, expected) {
		t.Errorf("Decode() = %v, want %v", result, expected)
	}
}

func TestUint16(t *testing.T) {
	input := []byte{0xCD, 0xFF, 0xFF}
	expected := uint16(65535)

	decoder := NewMessagePackDecoder(input)
	result, err := decoder.Decode()
	if err != nil {
		t.Errorf("Decode() error = %v", err)
	}
	if !isEqual(result, expected) {
		t.Errorf("Decode() = %v, want %v", result, expected)
	}
}

func TestUint32(t *testing.T) {
	input := []byte{0xCE, 0xFF, 0xFF, 0xFF, 0xFF}
	expected := uint32(4294967295)

	decoder := NewMessagePackDecoder(input)
	result, err := decoder.Decode()
	if err != nil {
		t.Errorf("Decode() error = %v", err)
	}
	if !isEqual(result, expected) {
		t.Errorf("Decode() = %v, want %v", result, expected)
	}
}

func TestUint64(t *testing.T) {
	input := []byte{0xCF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}
	expected := uint64(18446744073709551615)

	decoder := NewMessagePackDecoder(input)
	result, err := decoder.Decode()
	if err != nil {
		t.Errorf("Decode() error = %v", err)
	}
	if !isEqual(result, expected) {
		t.Errorf("Decode() = %v, want %v", result, expected)
	}
}

func TestInt8(t *testing.T) {
	input := []byte{0xD0, 0xFF}
	expected := int8(-1)

	decoder := NewMessagePackDecoder(input)
	result, err := decoder.Decode()
	if err != nil {
		t.Errorf("Decode() error = %v", err)
	}
	if !isEqual(result, expected) {
		t.Errorf("Decode() = %v, want %v", result, expected)
	}
}

func TestInt16(t *testing.T) {
	input := []byte{0xD1, 0xFF, 0xFF}
	expected := int16(-1)

	decoder := NewMessagePackDecoder(input)
	result, err := decoder.Decode()
	if err != nil {
		t.Errorf("Decode() error = %v", err)
	}
	if !isEqual(result, expected) {
		t.Errorf("Decode() = %v, want %v", result, expected)
	}
}

func TestInt32(t *testing.T) {
	input := []byte{0xD2, 0xFF, 0xFF, 0xFF, 0xFF}
	expected := int32(-1)

	decoder := NewMessagePackDecoder(input)
	result, err := decoder.Decode()
	if err != nil {
		t.Errorf("Decode() error = %v", err)
	}
	if !isEqual(result, expected) {
		t.Errorf("Decode() = %v, want %v", result, expected)
	}
}

func TestInt64(t *testing.T) {
	input := []byte{0xD3, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}
	expected := int64(-1)

	decoder := NewMessagePackDecoder(input)
	result, err := decoder.Decode()
	if err != nil {
		t.Errorf("Decode() error = %v", err)
	}
	if !isEqual(result, expected) {
		t.Errorf("Decode() = %v, want %v", result, expected)
	}
}

func TestStr8(t *testing.T) {
	input := []byte{0xD9, 0x21, 'f', 'o', 'o', 'b', 'a', 'r', 'b', 'a', 'z', 'q', 'u', 'x', 'q', 'u', 'u', 'x', 'o', 'l', 'o', 'r', 'e', 'm', 'i', 's', 'p', 's', 'u', 'm', 'd', 'o', 'l', 'o', 'r'}
	expected := "foobarbazquxquuxoloremispsumdolor"

	decoder := NewMessagePackDecoder(input)
	result, err := decoder.Decode()
	if err != nil {
		t.Errorf("Decode() error = %v", err)
	}
	if !isEqual(result, expected) {
		t.Errorf("Decode() = %v, want %v", result, expected)
	}
}

func TestArray16(t *testing.T) {
	input := []byte{0xDC, 0x00, 0x10, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10}
	expected := []interface{}{uint8(1), uint8(2), uint8(3), uint8(4), uint8(5), uint8(6), uint8(7), uint8(8), uint8(9), uint8(10), uint8(11), uint8(12), uint8(13), uint8(14), uint8(15), uint8(16)}

	decoder := NewMessagePackDecoder(input)
	result, err := decoder.Decode()
	if err != nil {
		t.Errorf("Decode() error = %v", err)
	}
	if !isEqual(result, expected) {
		t.Errorf("Decode() = %v, want %v", result, expected)
	}
}

func TestMap16(t *testing.T) {
	input := []byte{0xDE, 0x00, 0x10, 0xA4, 'k', 'e', 'y', '1', 0x01, 0xA4, 'k', 'e', 'y', '2', 0x02, 0xA4, 'k', 'e', 'y', '3', 0x03, 0xA4, 'k', 'e', 'y', '4', 0x04, 0xA4, 'k', 'e', 'y', '5', 0x05, 0xA4, 'k', 'e', 'y', '6', 0x06, 0xA4, 'k', 'e', 'y', '7', 0x07, 0xA4, 'k', 'e', 'y', '8', 0x08, 0xA4, 'k', 'e', 'y', '9', 0x09, 0xA5, 'k', 'e', 'y', '1', '0', 0x0A, 0xA5, 'k', 'e', 'y', '1', '1', 0x0B, 0xA5, 'k', 'e', 'y', '1', '2', 0x0C, 0xA5, 'k', 'e', 'y', '1', '3', 0x0D, 0xA5, 'k', 'e', 'y', '1', '4', 0x0E, 0xA5, 'k', 'e', 'y', '1', '5', 0x0F, 0xA5, 'k', 'e', 'y', '1', '6', 0x10}
	expected := map[string]interface{}{
		"key1": uint8(1), "key2": uint8(2), "key3": uint8(3), "key4": uint8(4),
		"key5": uint8(5), "key6": uint8(6), "key7": uint8(7), "key8": uint8(8),
		"key9": uint8(9), "key10": uint8(10), "key11": uint8(11), "key12": uint8(12),
		"key13": uint8(13), "key14": uint8(14), "key15": uint8(15), "key16": uint8(16),
	}

	decoder := NewMessagePackDecoder(input)
	result, err := decoder.Decode()
	if err != nil {
		t.Errorf("Decode() error = %v", err)
	}
	if !isEqual(result, expected) {
		t.Errorf("Decode() = %v, want %v", result, expected)
	}
}

func TestNegativeFixInt(t *testing.T) {
	tests := []struct {
		input    []byte
		expected interface{}
	}{
		{input: []byte{0xE0}, expected: int8(-32)},
		{input: []byte{0xFF}, expected: int8(-1)},
	}

	for _, test := range tests {
		decoder := NewMessagePackDecoder(test.input)
		result, err := decoder.Decode()
		if err != nil {
			t.Errorf("Decode() error = %v", err)
		}
		if !isEqual(result, test.expected) {
			t.Errorf("Decode() = %v, want %v", result, test.expected)
		}
	}
}

// Helper function to compare expected and actual values
func isEqual(a, b interface{}) bool {
	switch a := a.(type) {
	case []byte:
		bb, ok := b.([]byte)
		if !ok {
			return false
		}
		return bytes.Equal(a, bb)
	case []interface{}:
		bb, ok := b.([]interface{})
		if !ok {
			return false
		}
		if len(a) != len(bb) {
			return false
		}
		for i := range a {
			if !isEqual(a[i], bb[i]) {
				return false
			}
		}
		return true
	case map[string]interface{}:
		bb, ok := b.(map[string]interface{})
		if !ok {
			return false
		}
		if len(a) != len(bb) {
			return false
		}
		for key, value := range a {
			if !isEqual(value, bb[key]) {
				return false
			}
		}
		return true
	default:
		return a == b
	}
}
