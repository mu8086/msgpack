package msgpack

import (
	"bytes"
	"strings"
	"testing"
)

// TODO:
func Test_encode(t *testing.T) {
	var buf bytes.Buffer

	tests := []struct {
		name    string
		arg     interface{}
		wantErr error
	}{
		{name: "integer", arg: 1, wantErr: ErrUnsupportedType},
		{name: "string", arg: "", wantErr: nil},
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
		// TODO: {name: "want false, got true", args: args{false, []byte{0xC3}}, wantErr: false},
		{name: "want true, got true", args: args{true, []byte{0xC3}}, wantErr: false},
		// TODO: {name: "want true, got false", args: args{true, []byte{0xC2}}, wantErr: false},
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
				if b != prefixByte {
					t.Errorf("prefix not match, want: %v, got: %v\n", prefixByte, b)
				} else if err != nil {
					t.Errorf("buf.ReadByte() failed, err: %v\n", err)
				}
			}
		})
	}
}
