// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag_test

import (
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/gowarden/zflag"
)

func TestI64SValueImplementsGetter(t *testing.T) {
	f := zflag.NewFlagSet("test", zflag.ContinueOnError)
	f.Int64Slice("is", []int64{0, 1}, "Command separated list!")
	v := f.Lookup("is").Value

	if _, ok := v.(zflag.Getter); !ok {
		t.Fatalf("%T should implement the Getter interface", v)
	}
}

func TestInt64Slice(t *testing.T) {
	tests := []struct {
		name           string
		flagDefault    []int64
		input          []string
		expectedErr    string
		expectedValues []int64
		visitor        func(f *zflag.Flag)
	}{
		{
			name:           "no value passed",
			input:          []string{},
			flagDefault:    []int64{},
			expectedErr:    "",
			expectedValues: []int64{},
		},
		{
			name:        "empty value passed",
			input:       []string{""},
			flagDefault: []int64{},
			expectedErr: `invalid argument "" for "--i64s" flag: strconv.ParseInt: parsing "": invalid syntax`,
		},
		{
			name:        "invalid int64",
			input:       []string{"blabla"},
			flagDefault: []int64{},
			expectedErr: `invalid argument "blabla" for "--i64s" flag: strconv.ParseInt: parsing "blabla": invalid syntax`,
		},
		{
			name:        "no csv",
			input:       []string{"1,5"},
			flagDefault: []int64{},
			expectedErr: `invalid argument "1,5" for "--i64s" flag: strconv.ParseInt: parsing "1,5": invalid syntax`,
		},
		{
			name:           "empty defaults",
			input:          []string{"1", "5"},
			flagDefault:    []int64{},
			expectedValues: []int64{1, 5},
		},
		{
			name:           "with default values",
			input:          []string{"5", "1"},
			flagDefault:    []int64{1, 5},
			expectedValues: []int64{5, 1},
		},
		{
			name:           "trims input",
			input:          []string{"    1", "2    ", "   3  "},
			flagDefault:    []int64{},
			expectedValues: []int64{1, 2, 3},
		},
		{
			name:  "replace values",
			input: []string{"5", "1"},
			visitor: func(f *zflag.Flag) {
				if val, ok := f.Value.(zflag.SliceValue); ok {
					_ = val.Replace([]string{"3"})
				}
			},
			expectedValues: []int64{3},
		},
	}

	t.Parallel()
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			var i64s []int64
			f := zflag.NewFlagSet("test", zflag.ContinueOnError)
			f.SetOutput(ioutil.Discard)
			f.Int64SliceVar(&i64s, "i64s", test.flagDefault, "usage")
			err := f.Parse(repeatFlag("--i64s", test.input...))
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

			if !reflect.DeepEqual(test.expectedValues, i64s) {
				t.Fatalf("expected %v with type %T but got %v with type %T ", test.expectedValues, test.expectedValues, i64s, i64s)
			}

			int64Slice, err := f.GetInt64Slice("i64s")
			if err != nil {
				t.Fatal("got an error from GetInt64Slice():", err)
			}
			if !reflect.DeepEqual(test.expectedValues, int64Slice) {
				t.Fatalf("expected %v with type %T but got %v with type %T ", test.expectedValues, test.expectedValues, int64Slice, int64Slice)
			}

			int64SliceGet, err := f.Get("i64s")
			if err != nil {
				t.Fatal("got an error from Get():", err)
			}
			if !reflect.DeepEqual(int64SliceGet, int64Slice) {
				t.Fatalf("expected %v with type %T but got %v with type %T ", test.expectedValues, test.expectedValues, int64SliceGet, int64SliceGet)
			}
		})
	}
}
