// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag_test

import (
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/gowarden/zflag"
)

func TestStringToString(t *testing.T) {
	tests := []struct {
		name           string
		flagDefault    map[string]string
		input          []string
		expectedErr    string
		expectedValues map[string]string
		visitor        func(f *zflag.Flag)
	}{
		{
			name:           "no value passed",
			input:          []string{},
			flagDefault:    map[string]string{},
			expectedErr:    "",
			expectedValues: map[string]string{},
		},
		{
			name:        "empty value passed",
			input:       []string{""},
			flagDefault: map[string]string{},
			expectedErr: `invalid argument "" for "--s2s" flag:  must be formatted as key=value`,
		},
		{
			name:        "invalid string",
			input:       []string{"blabla"},
			flagDefault: map[string]string{},
			expectedErr: `invalid argument "blabla" for "--s2s" flag: blabla must be formatted as key=value`,
		},
		{
			name:           "no csv",
			input:          []string{"test=1,5"},
			flagDefault:    map[string]string{},
			expectedValues: map[string]string{"test": "1,5"},
		},
		{
			name:           "single key value pair per arg",
			input:          []string{"test=1=1"},
			flagDefault:    map[string]string{},
			expectedValues: map[string]string{"test": "1=1"},
		},
		{
			name:           "overrides multiple calls",
			input:          []string{"test=1", "test=5"},
			flagDefault:    map[string]string{},
			expectedValues: map[string]string{"test": "5"},
		},
		{
			name:           "empty defaults",
			input:          []string{"test=1", "test2=5"},
			flagDefault:    map[string]string{},
			expectedValues: map[string]string{"test": "1", "test2": "5"},
		},
		{
			name:           "overrides default values",
			input:          []string{"test=1", "test2=5"},
			flagDefault:    map[string]string{"test2": "1", "test": "5"},
			expectedValues: map[string]string{"test": "1", "test2": "5"},
		},
		{
			name:           "returns default values",
			input:          []string{},
			flagDefault:    map[string]string{"test2": "1", "test": "5"},
			expectedValues: map[string]string{"test2": "1", "test": "5"},
		},
		{
			name:           "keeps whitespace",
			input:          []string{"test1=asd   ", "test2=   value", "test3=    asd   ", "test4=multi\nline\narg\npassed\nin\n"},
			flagDefault:    map[string]string{},
			expectedValues: map[string]string{"test1": "asd   ", "test2": "   value", "test3": "    asd   ", "test4": "multi\nline\narg\npassed\nin\n"},
		},
	}

	t.Parallel()
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			var s2s map[string]string
			f := zflag.NewFlagSet("test", zflag.ContinueOnError)
			f.SetOutput(ioutil.Discard)
			f.StringToStringVar(&s2s, "s2s", test.flagDefault, "usage")
			err := f.Parse(repeatFlag("--s2s", test.input...))
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

			if !reflect.DeepEqual(test.expectedValues, s2s) {
				t.Fatalf("expected %[1]v with type %[1]T but got %[2]v with type %[2]T", test.expectedValues, s2s)
			}

			int16Slice, err := f.GetStringToString("s2s")
			if err != nil {
				t.Fatal("got an error from GetStringToString():", err)
			}
			if !reflect.DeepEqual(test.expectedValues, int16Slice) {
				t.Fatalf("expected %[1]v with type %[1]T but got %[2]v with type %[2]T", test.expectedValues, int16Slice)
			}

			int16SliceGet, err := f.Get("s2s")
			if err != nil {
				t.Fatal("got an error from Get():", err)
			}
			if !reflect.DeepEqual(test.expectedValues, int16SliceGet) {
				t.Fatalf("expected %[1]v with type %[1]T but got %[2]v with type %[2]T", test.expectedValues, int16SliceGet)
			}
		})
	}
}
