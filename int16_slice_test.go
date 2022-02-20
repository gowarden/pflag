// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag_test

import (
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/gowarden/zflag"
)

func TestI16SliceValueImplementsGetter(t *testing.T) {
	f := zflag.NewFlagSet("test", zflag.ContinueOnError)
	f.Int16Slice("is", []int16{0, 1}, "Command separated list!")
	v := f.Lookup("is").Value

	if _, ok := v.(zflag.Getter); !ok {
		t.Fatalf("%T should implement the Getter interface", v)
	}
}

func TestInt16Slice(t *testing.T) {
	tests := []struct {
		name           string
		flagDefault    []int16
		input          []string
		expectedErr    string
		expectedValues []int16
		visitor        func(f *zflag.Flag)
	}{
		{
			name:           "no value passed",
			input:          []string{},
			flagDefault:    []int16{},
			expectedErr:    "",
			expectedValues: []int16{},
		},
		{
			name:        "empty value passed",
			input:       []string{""},
			flagDefault: []int16{},
			expectedErr: `invalid argument "" for "--i16s" flag: strconv.ParseInt: parsing "": invalid syntax`,
		},
		{
			name:        "invalid int16",
			input:       []string{"blabla"},
			flagDefault: []int16{},
			expectedErr: `invalid argument "blabla" for "--i16s" flag: strconv.ParseInt: parsing "blabla": invalid syntax`,
		},
		{
			name:        "no csv",
			input:       []string{"1,5"},
			flagDefault: []int16{},
			expectedErr: `invalid argument "1,5" for "--i16s" flag: strconv.ParseInt: parsing "1,5": invalid syntax`,
		},
		{
			name:           "empty defaults",
			input:          []string{"1", "5"},
			flagDefault:    []int16{},
			expectedValues: []int16{1, 5},
		},
		{
			name:           "with default values",
			input:          []string{"5", "1"},
			flagDefault:    []int16{1, 5},
			expectedValues: []int16{5, 1},
		},
		{
			name:           "trims input",
			input:          []string{"    1", "2    ", "   3  "},
			flagDefault:    []int16{},
			expectedValues: []int16{1, 2, 3},
		},
		{
			name:  "replace values",
			input: []string{"5", "1"},
			visitor: func(f *zflag.Flag) {
				if val, ok := f.Value.(zflag.SliceValue); ok {
					_ = val.Replace([]string{"3"})
				}
			},
			expectedValues: []int16{3},
		},
	}

	t.Parallel()
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			var i16s []int16
			f := zflag.NewFlagSet("test", zflag.ContinueOnError)
			f.SetOutput(ioutil.Discard)
			f.Int16SliceVar(&i16s, "i16s", test.flagDefault, "usage")
			err := f.Parse(repeatFlag("--i16s", test.input...))
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

			if test.visitor != nil {
				f.VisitAll(test.visitor)
			}

			if !reflect.DeepEqual(test.expectedValues, i16s) {
				t.Fatalf("expected %v with type %T but got %v with type %T ", test.expectedValues, test.expectedValues, i16s, i16s)
			}

			int16Slice, err := f.GetInt16Slice("i16s")
			if err != nil {
				t.Fatal("got an error from GetInt16Slice():", err)
			}
			if !reflect.DeepEqual(test.expectedValues, int16Slice) {
				t.Fatalf("expected %v with type %T but got %v with type %T ", test.expectedValues, test.expectedValues, int16Slice, int16Slice)
			}

			int16SliceGet, err := f.Get("i16s")
			if err != nil {
				t.Fatal("got an error from Get():", err)
			}
			if !reflect.DeepEqual(int16SliceGet, int16Slice) {
				t.Fatalf("expected %v with type %T but got %v with type %T ", test.expectedValues, test.expectedValues, int16SliceGet, int16SliceGet)
			}
		})
	}
}
