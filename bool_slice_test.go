// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag_test

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/gowarden/zflag"
)

func TestBoolSlice(t *testing.T) {
	tests := []struct {
		name           string
		flagDefault    []bool
		input          []string
		expectedErr    string
		expectedValues []bool
		visitor        func(f *zflag.Flag)
	}{
		{
			name:           "no value passed",
			input:          []string{},
			flagDefault:    []bool{},
			expectedErr:    "",
			expectedValues: []bool{},
		},
		{
			name:        "empty value passed",
			input:       []string{""},
			flagDefault: []bool{},
			expectedErr: `invalid argument "" for "--bs" flag: strconv.ParseBool: parsing "": invalid syntax`,
		},
		{
			name:        "invalid bool",
			input:       []string{"blabla"},
			flagDefault: []bool{},
			expectedErr: `invalid argument "blabla" for "--bs" flag: strconv.ParseBool: parsing "blabla": invalid syntax`,
		},
		{
			name:        "no csv",
			input:       []string{"true,false"},
			flagDefault: []bool{},
			expectedErr: `invalid argument "true,false" for "--bs" flag: strconv.ParseBool: parsing "true,false": invalid syntax`,
		},
		{
			name:           "empty value passed",
			input:          []string{"true", "false"},
			flagDefault:    []bool{},
			expectedValues: []bool{true, false},
		},
		{
			name:           "with default values",
			input:          []string{"false", "true"},
			flagDefault:    []bool{true, false},
			expectedValues: []bool{false, true},
		},
		{
			name:  "replace values",
			input: []string{"true", "false"},
			visitor: func(f *zflag.Flag) {
				if val, ok := f.Value.(zflag.SliceValue); ok {
					_ = val.Replace([]string{"false"})
				}
			},
			expectedValues: []bool{false},
		},
		{
			name:           "trims input",
			input:          []string{" true ", " false "},
			expectedValues: []bool{true, false},
		},
		{
			name:           "all valid bool values",
			input:          []string{"true", "false", "1", "0", "t", "f", "TRUE", "FALSE", "1", "0", "T", "F", "True", "False"},
			expectedValues: []bool{true, false, true, false, true, false, true, false, true, false, true, false, true, false},
		},
	}

	t.Parallel()
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			var bs []bool
			f := zflag.NewFlagSet("test", zflag.ContinueOnError)
			f.SetOutput(ioutil.Discard)
			f.BoolSliceVar(&bs, "bs", test.flagDefault, "usage")
			err := f.Parse(repeatFlag("--bs", test.input...))
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

			if !reflect.DeepEqual(test.expectedValues, bs) {
				t.Fatalf("expected %[1]v with type %[1]T but got %[2]v with type %[2]T", test.expectedValues, bs)
			}

			boolSlice, err := f.GetBoolSlice("bs")
			if err != nil {
				t.Fatal("got an error from GetBoolSlice():", err)
			}
			if !reflect.DeepEqual(test.expectedValues, boolSlice) {
				t.Fatalf("expected %[1]v with type %[1]T but got %[2]v with type %[2]T", test.expectedValues, boolSlice)
			}

			boolSliceGet, err := f.Get("bs")
			if err != nil {
				t.Fatal("got an error from Get():", err)
			}
			if !reflect.DeepEqual(boolSliceGet, boolSlice) {
				t.Fatalf("expected %[1]v with type %[1]T but got %[2]v with type %[2]T", test.expectedValues, boolSliceGet)
			}
		})
	}
}

func repeatFlag(flag string, values ...string) (res []string) {
	res = make([]string, 0, len(values))
	for _, val := range values {
		res = append(res, fmt.Sprintf("%s=%s", flag, val))
	}

	return
}
