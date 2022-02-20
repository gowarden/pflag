// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag_test

import (
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/gowarden/zflag"
)

func TestF64SValueImplementsGetter(t *testing.T) {
	f := zflag.NewFlagSet("test", zflag.ContinueOnError)
	f.Float64Slice("f64s", []float64{0.0, 1.0}, "Command separated list!")
	v := f.Lookup("f64s").Value

	if _, ok := v.(zflag.Getter); !ok {
		t.Fatalf("%T should implement the Getter interface", v)
	}
}

func TestFloat64Slice(t *testing.T) {
	tests := []struct {
		name           string
		flagDefault    []float64
		input          []string
		expectedErr    string
		expectedValues []float64
		visitor        func(f *zflag.Flag)
	}{
		{
			name:           "no value passed",
			input:          []string{},
			flagDefault:    []float64{},
			expectedErr:    "",
			expectedValues: []float64{},
		},
		{
			name:        "empty value passed",
			input:       []string{""},
			flagDefault: []float64{},
			expectedErr: `invalid argument "" for "--f64s" flag: strconv.ParseFloat: parsing "": invalid syntax`,
		},
		{
			name:        "invalid float64",
			input:       []string{"blabla"},
			flagDefault: []float64{},
			expectedErr: `invalid argument "blabla" for "--f64s" flag: strconv.ParseFloat: parsing "blabla": invalid syntax`,
		},
		{
			name:        "no csv",
			input:       []string{"1.1,1.5"},
			flagDefault: []float64{},
			expectedErr: `invalid argument "1.1,1.5" for "--f64s" flag: strconv.ParseFloat: parsing "1.1,1.5": invalid syntax`,
		},
		{
			name:           "empty value passed",
			input:          []string{"1.5", "1.1"},
			flagDefault:    []float64{},
			expectedValues: []float64{1.5, 1.1},
		},
		{
			name:           "with default values",
			input:          []string{"1.5", "1.1"},
			flagDefault:    []float64{1.1, 1.5},
			expectedValues: []float64{1.5, 1.1},
		},
		{
			name:           "trims input",
			input:          []string{"    1.5", "1.1    ", "   1.1  "},
			flagDefault:    []float64{},
			expectedValues: []float64{1.5, 1.1, 1.1},
		},
		{
			name:  "replace values",
			input: []string{"1.5", "1.1"},
			visitor: func(f *zflag.Flag) {
				if val, ok := f.Value.(zflag.SliceValue); ok {
					_ = val.Replace([]string{"1.3"})
				}
			},
			expectedValues: []float64{1.3},
		},
	}

	t.Parallel()
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			var f64s []float64
			f := zflag.NewFlagSet("test", zflag.ContinueOnError)
			f.SetOutput(ioutil.Discard)
			f.Float64SliceVar(&f64s, "f64s", test.flagDefault, "usage")
			err := f.Parse(repeatFlag("--f64s", test.input...))
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

			if !reflect.DeepEqual(test.expectedValues, f64s) {
				t.Fatalf("expected %v with type %T but got %v with type %T ", test.expectedValues, test.expectedValues, f64s, f64s)
			}

			float64Slice, err := f.GetFloat64Slice("f64s")
			if err != nil {
				t.Fatal("got an error from GetFloat64Slice():", err)
			}
			if !reflect.DeepEqual(test.expectedValues, float64Slice) {
				t.Fatalf("expected %v with type %T but got %v with type %T ", test.expectedValues, test.expectedValues, float64Slice, float64Slice)
			}

			float64SliceGet, err := f.Get("f64s")
			if err != nil {
				t.Fatal("got an error from Get():", err)
			}
			if !reflect.DeepEqual(float64SliceGet, float64Slice) {
				t.Fatalf("expected %v with type %T but got %v with type %T ", test.expectedValues, test.expectedValues, float64SliceGet, float64SliceGet)
			}
		})
	}
}
