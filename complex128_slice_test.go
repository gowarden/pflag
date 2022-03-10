// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build go1.15
// +build go1.15

package zflag_test

import (
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/gowarden/zflag"
)

func TestC128Slice(t *testing.T) {
	tests := []struct {
		name           string
		flagDefault    []complex128
		input          []string
		expectedErr    string
		expectedValues []complex128
		visitor        func(f *zflag.Flag)
	}{
		{
			name:           "no value passed",
			input:          []string{},
			flagDefault:    []complex128{},
			expectedErr:    "",
			expectedValues: []complex128{},
		},
		{
			name:        "empty value passed",
			input:       []string{""},
			flagDefault: []complex128{},
			expectedErr: `invalid argument "" for "--c128s" flag: strconv.ParseComplex: parsing "": invalid syntax`,
		},
		{
			name:        "invalid c128s",
			input:       []string{"blabla"},
			flagDefault: []complex128{},
			expectedErr: `invalid argument "blabla" for "--c128s" flag: strconv.ParseComplex: parsing "blabla": invalid syntax`,
		},
		{
			name:        "no csv",
			input:       []string{"1.0,2.0"},
			flagDefault: []complex128{},
			expectedErr: `invalid argument "1.0,2.0" for "--c128s" flag: strconv.ParseComplex: parsing "1.0,2.0": invalid syntax`,
		},
		{
			name:           "multiple values passed",
			input:          []string{"1.0", "2.0"},
			flagDefault:    []complex128{},
			expectedValues: []complex128{1.0, 2.0},
		},
		{
			name:           "with default values",
			input:          []string{"1.0", "2.0"},
			flagDefault:    []complex128{2.0, 1.0},
			expectedValues: []complex128{1.0, 2.0},
		},
		{
			name:  "replace values",
			input: []string{"1.0", "2.0"},
			visitor: func(f *zflag.Flag) {
				if val, ok := f.Value.(zflag.SliceValue); ok {
					_ = val.Replace([]string{"0+2i"})
				}
			},
			expectedValues: []complex128{complex(0, 2)},
		},
		{
			name:           "valid c128s",
			input:          []string{"1.0", "2.0", "3.0", "0+2i", "1", "2i", "2.5+3.1i"},
			expectedValues: []complex128{1.0, 2.0, 3.0, complex(0, 2), complex(1, 0), complex(0, 2), complex(2.5, 3.1)},
		},
		{
			name:           "trims input",
			input:          []string{" 1.0 ", "   2.0", "3.0   ", "  0+2i", "1"},
			expectedValues: []complex128{1.0, 2.0, 3.0, complex(0, 2), complex(1, 0)},
		},
	}

	t.Parallel()
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			var c128s []complex128
			f := zflag.NewFlagSet("test", zflag.ContinueOnError)
			f.SetOutput(ioutil.Discard)
			f.Complex128SliceVar(&c128s, "c128s", test.flagDefault, "usage")
			err := f.Parse(repeatFlag("--c128s", test.input...))
			if test.expectedErr != "" {
				if err == nil {
					t.Fatalf("expected an error; got none")
				}
				if test.expectedErr != "" && err.Error() != test.expectedErr {
					t.Fatalf("expected error to equal %q, but was: %s", test.expectedErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("expected no error; got %q", err)
			}

			if test.visitor != nil {
				f.VisitAll(test.visitor)
			}

			if !reflect.DeepEqual(test.expectedValues, c128s) {
				t.Fatalf("expected %[1]v with type %[1]T but got %[2]v with type %[2]T", test.expectedValues, c128s)
			}

			getC128s, err := f.GetComplex128Slice("c128s")
			if err != nil {
				t.Fatal("got an error from GetComplex128Slice():", err)
			}
			if !reflect.DeepEqual(test.expectedValues, getC128s) {
				t.Fatalf("expected %[1]v with type %[1]T but got %[2]v with type %[2]T", test.expectedValues, getC128s)
			}

			getC128sGet, err := f.Get("c128s")
			if err != nil {
				t.Fatal("got an error from Get():", err)
			}
			if !reflect.DeepEqual(getC128sGet, c128s) {
				t.Fatalf("expected %[1]v with type %[1]T but got %[2]v with type %[2]T", test.expectedValues, getC128sGet)
			}
		})
	}
}
