// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/gowarden/zflag"
)

func TestBytesValueImplementsGetter(t *testing.T) {
	f := zflag.NewFlagSet("test", zflag.ContinueOnError)
	f.BytesBase64("b64", []byte{}, "b64")
	f.BytesHex("hex", []byte{}, "hex")

	v := f.Lookup("b64").Value

	if _, ok := v.(zflag.Getter); !ok {
		t.Fatalf("%T should implement the Getter interface", v)
	}

	v2 := f.Lookup("b64").Value

	if _, ok := v2.(zflag.Getter); !ok {
		t.Fatalf("%T should implement the Getter interface", v2)
	}
}

func TestBytesHex(t *testing.T) {
	tests := []struct {
		name          string
		flagDefault   []byte
		input         []string
		expectedErr   string
		expectedValue string
	}{
		{
			name:          "no value passed",
			input:         []string{},
			flagDefault:   []byte{},
			expectedErr:   "",
			expectedValue: "",
		},
		{
			name:          "empty value passed",
			input:         []string{""},
			flagDefault:   []byte{},
			expectedErr:   "",
			expectedValue: "",
		},
		{
			name:        "invalid byte hex short string",
			input:       []string{"0"},
			flagDefault: []byte{},
			expectedErr: `invalid argument "0" for "--bytes" flag: encoding/hex: odd length hex string`,
		},
		{
			name:        "invalid byte hex odd-length string",
			input:       []string{"000"},
			flagDefault: []byte{},
			expectedErr: `invalid argument "000" for "--bytes" flag: encoding/hex: odd length hex string`,
		},
		{
			name:        "invalid byte hex non-hex char",
			input:       []string{"qq"},
			flagDefault: []byte{},
			expectedErr: `invalid argument "qq" for "--bytes" flag: encoding/hex: invalid byte: U+0071 'q'`,
		},
		{
			name:        "no csv",
			input:       []string{"0101,0101"},
			flagDefault: []byte{},
			expectedErr: `invalid argument "0101,0101" for "--bytes" flag: encoding/hex: invalid byte: U+002C ','`,
		},
		{
			name:          "repeated value gets the last call",
			input:         []string{"01", "0101"},
			flagDefault:   []byte{},
			expectedValue: "0101",
		},
		{
			name:          "default values get overwritten",
			input:         []string{"0101"},
			flagDefault:   []byte("01"),
			expectedValue: "0101",
		},
		{
			name:          "trims input",
			input:         []string{" 01 "},
			expectedValue: "01",
		},
		{
			name:          "test valid",
			input:         []string{"01"},
			expectedValue: "01",
		},
		{
			name:          "test valid",
			input:         []string{"0101"},
			expectedValue: "0101",
		},
		{
			name:          "test valid",
			input:         []string{"1234567890abcdef"},
			expectedValue: "1234567890ABCDEF",
		},
		{
			name:          "test valid",
			input:         []string{"1234567890ABCDEF"},
			expectedValue: "1234567890ABCDEF",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var bytes []byte
			f := zflag.NewFlagSet("test", zflag.ContinueOnError)
			f.BytesHexVar(&bytes, "bytes", test.flagDefault, "usage")
			err := f.Parse(repeatFlag("--bytes", test.input...))
			if test.expectedErr != "" {
				if err == nil {
					t.Fatalf("expected an error; got none")
				}
				if test.expectedErr != "" && err.Error() != test.expectedErr {
					t.Fatalf("expected error to eqaul %q, but was: %s", test.expectedErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("expected no error; got %q", err)
			}

			if fmt.Sprintf("%X", bytes) != test.expectedValue {
				t.Fatalf("expected %v with type %T but got %v with type %T ", test.expectedValue, test.expectedValue, bytes, bytes)
			}

			bytesHex, err := f.GetBytesHex("bytes")
			if err != nil {
				t.Fatal("got an error from GetBytesHex():", err)
			}
			if fmt.Sprintf("%X", bytesHex) != test.expectedValue {
				t.Fatalf("expected %v with type %T but got %v with type %T ", bytesHex, bytesHex, test.expectedValue, test.expectedValue)
			}

			bytesHexGet, err := f.Get("bytes")
			if err != nil {
				t.Fatal("got an error from Get():", err)
			}
			if fmt.Sprintf("%X", bytesHexGet) != test.expectedValue {
				t.Fatalf("expected %v with type %T but got %v with type %T ", bytesHex, bytesHex, bytesHexGet, bytesHexGet)
			}
		})
	}
}

func TestBytesB64(t *testing.T) {
	tests := []struct {
		name          string
		flagDefault   []byte
		input         []string
		expectedErr   string
		expectedValue []byte
	}{
		{
			name:          "no value passed",
			input:         []string{},
			flagDefault:   []byte{},
			expectedValue: []byte{},
		},
		{
			name:          "empty value passed",
			input:         []string{""},
			flagDefault:   []byte{},
			expectedValue: []byte{},
		},
		{
			name:        "invalid byte base64 short string",
			input:       []string{"A"},
			flagDefault: []byte{},
			expectedErr: `invalid argument "A" for "--bytes" flag: illegal base64 data at input byte 0`,
		},
		{
			name:        "invalid byte hex non-hex char",
			input:       []string{"Aï=="},
			flagDefault: []byte{},
			expectedErr: `invalid argument "Aï==" for "--bytes" flag: illegal base64 data at input byte 1`,
		},
		{
			name:        "no csv",
			input:       []string{"AQ==,AQ=="},
			flagDefault: []byte{},
			expectedErr: `invalid argument "AQ==,AQ==" for "--bytes" flag: illegal base64 data at input byte 4`,
		},
		{
			name:          "default value gets returned when none passed",
			input:         []string{},
			flagDefault:   []byte("bye"),
			expectedValue: []byte("bye"),
		},
		{
			name:          "repeated value gets the last call",
			input:         []string{"aGk=", "Ynll"},
			flagDefault:   []byte{},
			expectedValue: []byte("bye"),
		},
		{
			name:          "default values get overwritten",
			input:         []string{"Ynll"},
			flagDefault:   []byte("aGk="),
			expectedValue: []byte("bye"),
		},
		{
			name:          "trims input",
			input:         []string{" Ynll "},
			expectedValue: []byte("bye"),
		},
		{
			name:          "test valid",
			input:         []string{"Ynll"},
			expectedValue: []byte("bye"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var bytes []byte
			f := zflag.NewFlagSet("test", zflag.ContinueOnError)
			f.BytesBase64Var(&bytes, "bytes", test.flagDefault, "usage")
			err := f.Parse(repeatFlag("--bytes", test.input...))
			if test.expectedErr != "" {
				if err == nil {
					t.Fatalf("expected an error; got none")
				}
				if test.expectedErr != "" && err.Error() != test.expectedErr {
					t.Fatalf("expected error to eqaul %q, but was: %s", test.expectedErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("expected no error; got %q", err)
			}

			if !reflect.DeepEqual(bytes, test.expectedValue) {
				t.Fatalf("expected %v with type %T but got %v with type %T ", test.expectedValue, test.expectedValue, bytes, bytes)
			}

			bytesB64, err := f.GetBytesBase64("bytes")
			if err != nil {
				t.Fatal("got an error from GetBytesHex():", err)
			}
			if !reflect.DeepEqual(bytesB64, test.expectedValue) {
				t.Fatalf("expected %v with type %T but got %v with type %T ", test.expectedValue, test.expectedValue, bytesB64, bytesB64)
			}

			bytesB64Get, err := f.Get("bytes")
			if err != nil {
				t.Fatal("got an error from Get():", err)
			}
			if !reflect.DeepEqual(bytesB64Get, test.expectedValue) {
				t.Fatalf("expected %v with type %T but got %v with type %T ", test.expectedValue, test.expectedValue, bytesB64Get, bytesB64Get)
			}
		})
	}
}
