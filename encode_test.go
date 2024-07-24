package msgpack

import (
	"bytes"
	"strconv"
	"strings"
	"testing"
)

func Test_encode(t *testing.T) {
	var buf bytes.Buffer

	tests := []struct {
		name    string
		arg     interface{}
		wantErr error
	}{
		{name: "integer", arg: 1, wantErr: nil},
		{name: "string", arg: "", wantErr: nil},
		{name: "bool true", arg: true, wantErr: nil},
		{name: "bool false", arg: false, wantErr: nil},
		{name: "nil", arg: nil, wantErr: nil},
		{name: "float", arg: 1.23, wantErr: nil},
		{name: "array", arg: []interface{}{1, "test", true}, wantErr: nil},
		{name: "map", arg: map[string]interface{}{"key": "value", "number": 42}, wantErr: nil},
		{name: "unsupported type", arg: struct{}{}, wantErr: ErrUnsupportedType},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()

			if err := encode(&buf, tt.arg); err != tt.wantErr {
				t.Errorf("encode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_encodeBool(t *testing.T) {
	var buf bytes.Buffer

	type args struct {
		value   bool
		encoded []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "want false, got false", args: args{false, []byte{0xC2}}, wantErr: false},
		{name: "want true, got true", args: args{true, []byte{0xC3}}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()

			if err := encodeBool(&buf, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("encodeBool() error = %v, wantErr %v", err, tt.wantErr)
			}

			b, err := buf.ReadByte()
			if b != tt.args.encoded[0] {
				t.Errorf("encoded not match, want: %v, got: %v\n", tt.args.encoded[0], b)
			} else if err != nil {
				t.Errorf("buf.ReadByte() failed, err: %v\n", err)
			}
		})
	}
}

func Test_encodeString(t *testing.T) {
	var buf bytes.Buffer

	type args struct {
		value         string
		encodedPrefix []byte
	}

	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{name: "fixstr, size 0", args: args{"", []byte{0xA0}}, wantErr: nil},
		{name: "fixstr, size 1", args: args{"A", []byte{0xA1}}, wantErr: nil},
		{name: "fixstr, size 31", args: args{strings.Repeat("A", 31), []byte{0xBF}}, wantErr: nil},
		{name: "str 8, size 32", args: args{strings.Repeat("A", 32), []byte{0xD9, 0x20}}, wantErr: nil},
		{name: "str 8, size 255", args: args{strings.Repeat("A", 255), []byte{0xD9, 0xFF}}, wantErr: nil},
		{name: "str 16, size 256", args: args{strings.Repeat("A", 256), []byte{0xDA, 0x01, 0x00}}, wantErr: nil},
		{name: "str 16, size 4581", args: args{strings.Repeat("A", 4581), []byte{0xDA, 0x11, 0xE5}}, wantErr: nil},
		{name: "str 16, size 65535", args: args{strings.Repeat("A", 65535), []byte{0xDA, 0xFF, 0xFF}}, wantErr: nil},
		{name: "str 32, size 65536", args: args{strings.Repeat("A", 65536), []byte{0xDB, 0x00, 0x01, 0x00, 0x00}}, wantErr: nil},
		{name: "str 32, size 71251289", args: args{strings.Repeat("A", 71251289), []byte{0xDB, 0x04, 0x3F, 0x35, 0x59}}, wantErr: nil},
		{name: "str 32, size 4294967295", args: args{strings.Repeat("A", 4294967295), []byte{0xDB, 0xFF, 0xFF, 0xFF, 0xFF}}, wantErr: nil},
		{name: "StringTooLong, size 4294967296", args: args{strings.Repeat("A", 4294967296), []byte{}}, wantErr: ErrStringTooLong},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()

			if err := encodeString(&buf, tt.args.value); err != tt.wantErr {
				t.Errorf("encodeString() error = %v, wantErr %v", err, tt.wantErr)
			}

			for _, prefixByte := range tt.args.encodedPrefix {
				b, err := buf.ReadByte()
				if err != nil {
					t.Errorf("buf.ReadByte() failed, err: %v\n", err)
					break
				}
				if b != prefixByte {
					t.Errorf("encoded prefix not match, want: %v, got: %v\n", prefixByte, b)
				}
			}
		})
	}
}

func Test_encodeFloat(t *testing.T) {
	var buf bytes.Buffer

	tests := []struct {
		name    string
		value   float64
		encoded []byte
		wantErr bool
	}{
		{name: "float32", value: float64(float32(1.23)), encoded: []byte{0xCA, 0x3F, 0x9D, 0x70, 0xA4}, wantErr: false},
		{name: "float64", value: 1.23, encoded: []byte{0xCB, 0x3F, 0xF3, 0xAE, 0x14, 0x7A, 0xE1, 0x47, 0xAE}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()

			if err := encodeFloat(&buf, tt.value); (err != nil) != tt.wantErr {
				t.Errorf("encodeFloat() error = %v, wantErr %v", err, tt.wantErr)
			}

			for _, b := range tt.encoded {
				got, err := buf.ReadByte()
				if err != nil {
					t.Errorf("buf.ReadByte() failed, err: %v\n", err)
					break
				}
				if got != b {
					t.Errorf("encoded not match, want: %v, got: %v\n", b, got)
				}
			}
		})
	}
}

func Test_encodeArray(t *testing.T) {
	var buf bytes.Buffer

	tests := []struct {
		name    string
		value   []interface{}
		encoded []byte
		wantErr bool
	}{
		{name: "empty array", value: []interface{}{}, encoded: []byte{0x90}, wantErr: false},
		{name: "fixarray", value: []interface{}{1, "test", true}, encoded: []byte{0x93, 0x01, 0xA4, 't', 'e', 's', 't', 0xC3}, wantErr: false},
		{name: "array 16", value: make([]interface{}, 20), encoded: []byte{0xDC, 0x00, 0x14}, wantErr: false},
		{name: "array 32", value: make([]interface{}, 70000), encoded: []byte{0xDD, 0x00, 0x01, 0x11, 0x70}, wantErr: false},
		// {name: "array too long", value: make([]interface{}, 0x100000000), encoded: nil, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()

			if err := encodeArray(&buf, tt.value); (err != nil) != tt.wantErr {
				t.Errorf("encodeArray() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.encoded != nil {
				for _, b := range tt.encoded {
					got, err := buf.ReadByte()
					if err != nil {
						t.Errorf("buf.ReadByte() failed, err: %v\n", err)
						break
					}
					if got != b {
						t.Errorf("encoded not match, want: %v, got: %v\n", b, got)
					}
				}
			}
		})
	}
}

func Test_encodeMap(t *testing.T) {
	var buf bytes.Buffer

	mapWithSize := func(size int) map[string]interface{} {
		m := make(map[string]interface{})
		for i := 0; i < size; i++ {
			m[strconv.Itoa(i)] = 1
		}
		return m
	}

	map20 := mapWithSize(20)
	map70000 := mapWithSize(70000)

	tests := []struct {
		name    string
		value   map[string]interface{}
		encoded []byte
		wantErr bool
	}{
		{name: "empty map", value: map[string]interface{}{}, encoded: []byte{0x80}, wantErr: false},
		{name: "fixmap", value: map[string]interface{}{"key": "value"}, encoded: []byte{0x81, 0xA3, 'k', 'e', 'y', 0xA5, 'v', 'a', 'l', 'u', 'e'}, wantErr: false},
		{name: "map 16", value: map20, encoded: []byte{0xDE, 0x00, 0x14}, wantErr: false},
		{name: "map 32", value: map70000, encoded: []byte{0xDF, 0x00, 0x01, 0x11, 0x70}, wantErr: false},
		// {name: "map too long", value: make(map[string]interface{}, 0x100000000), encoded: nil, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()

			if err := encodeMap(&buf, tt.value); (err != nil) != tt.wantErr {
				t.Errorf("encodeMap() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.encoded != nil {
				for _, b := range tt.encoded {
					got, err := buf.ReadByte()
					if err != nil {
						t.Errorf("buf.ReadByte() failed, err: %v\n", err)
						break
					}
					if got != b {
						t.Errorf("encoded not match, want: %v, got: %v\n", b, got)
					}
				}
			}
		})
	}
}
