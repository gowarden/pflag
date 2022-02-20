// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag_test

import (
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/gowarden/zflag"
)

func TestS2IValueImplementsGetter(t *testing.T) {
	f := zflag.NewFlagSet("test", zflag.ContinueOnError)
	f.StringToInt("s2i", map[string]int{}, "Command separated ls2it!")
	v := f.Lookup("s2i").Value

	if _, ok := v.(zflag.Getter); !ok {
		t.Fatalf("%T should implement the Getter interface", v)
	}
}

func TestStringToInt(t *testing.T) {
	tests := []struct {
		name           string
		flagDefault    map[string]int
		input          []string
		expectedErr    string
		expectedValues map[string]int
		visitor        func(f *zflag.Flag)
	}{
		{
			name:           "no value passed",
			input:          []string{},
			flagDefault:    map[string]int{},
			expectedErr:    "",
			expectedValues: map[string]int{},
		},
		{
			name:        "empty value passed",
			input:       []string{""},
			flagDefault: map[string]int{},
			expectedErr: `invalid argument "" for "--s2i" flag:  must be formatted as key=value`,
		},
		{
			name:        "invalid int",
			input:       []string{"blabla"},
			flagDefault: map[string]int{},
			expectedErr: `invalid argument "blabla" for "--s2i" flag: blabla must be formatted as key=value`,
		},
		{
			name:        "no csv",
			input:       []string{"test=1,5"},
			flagDefault: map[string]int{},
			expectedErr: `invalid argument "test=1,5" for "--s2i" flag: strconv.Atoi: parsing "1,5": invalid syntax`,
		},
		{
			name:        "single key value pair per arg",
			input:       []string{"test=1=1"},
			flagDefault: map[string]int{},
			expectedErr: `invalid argument "test=1=1" for "--s2i" flag: strconv.Atoi: parsing "1=1": invalid syntax`,
		},
		{
			name:           "overrides multiple calls",
			input:          []string{"test=1", "test=5"},
			flagDefault:    map[string]int{},
			expectedValues: map[string]int{"test": 5},
		},
		{
			name:           "empty defaults",
			input:          []string{"test=1", "test2=5"},
			flagDefault:    map[string]int{},
			expectedValues: map[string]int{"test": 1, "test2": 5},
		},
		{
			name:           "overrides default values",
			input:          []string{"test=1", "test2=5"},
			flagDefault:    map[string]int{"test2": 1, "test": 5},
			expectedValues: map[string]int{"test": 1, "test2": 5},
		},
		{
			name:           "returns default values",
			input:          []string{},
			flagDefault:    map[string]int{"test2": 1, "test": 5},
			expectedValues: map[string]int{"test2": 1, "test": 5},
		},
		{
			name:           "trims input",
			input:          []string{"test=    1", "test2=5     ", "test3=     9     "},
			flagDefault:    map[string]int{},
			expectedValues: map[string]int{"test": 1, "test2": 5, "test3": 9},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var s2i map[string]int
			f := zflag.NewFlagSet("test", zflag.ContinueOnError)
			f.SetOutput(ioutil.Discard)
			f.StringToIntVar(&s2i, "s2i", test.flagDefault, "usage")
			err := f.Parse(repeatFlag("--s2i", test.input...))
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

			if !reflect.DeepEqual(test.expectedValues, s2i) {
				t.Fatalf("expected %v with type %T but got %v with type %T ", test.expectedValues, test.expectedValues, s2i, s2i)
			}

			int16Slice, err := f.GetStringToInt("s2i")
			if err != nil {
				t.Fatal("got an error from GetStringToInt():", err)
			}
			if !reflect.DeepEqual(test.expectedValues, int16Slice) {
				t.Fatalf("expected %v with type %T but got %v with type %T ", test.expectedValues, test.expectedValues, int16Slice, int16Slice)
			}

			int16SliceGet, err := f.Get("s2i")
			if err != nil {
				t.Fatal("got an error from Get():", err)
			}
			if !reflect.DeepEqual(int16SliceGet, int16Slice) {
				t.Fatalf("expected %v with type %T but got %v with type %T ", test.expectedValues, test.expectedValues, int16SliceGet, int16SliceGet)
			}
		})
	}
}
