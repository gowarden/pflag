// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag_test

import (
	"io/ioutil"
	"reflect"
	"testing"
	"time"

	"github.com/gowarden/zflag"
)

func TestDurationSlice(t *testing.T) {
	tests := []struct {
		name           string
		flagDefault    []time.Duration
		input          []string
		expectedErr    string
		expectedValues []time.Duration
		visitor        func(f *zflag.Flag)
	}{
		{
			name:           "no value passed",
			input:          []string{},
			flagDefault:    []time.Duration{},
			expectedErr:    "",
			expectedValues: []time.Duration{},
		},
		{
			name:        "empty value passed",
			input:       []string{""},
			flagDefault: []time.Duration{},
			expectedErr: `invalid argument "" for "--ds" flag: time: invalid duration ""`,
		},
		{
			name:        "invalid unit",
			input:       []string{"2q"},
			flagDefault: []time.Duration{},
			expectedErr: `invalid argument "2q" for "--ds" flag: time: unknown unit "q" in duration "2q"`,
		},
		{
			name:        "invalid duration",
			input:       []string{"blabla"},
			flagDefault: []time.Duration{},
			expectedErr: `invalid argument "blabla" for "--ds" flag: time: invalid duration "blabla"`,
		},
		{
			name:        "no csv",
			input:       []string{"2m,2h"},
			flagDefault: []time.Duration{},
			expectedErr: `invalid argument "2m,2h" for "--ds" flag: time: unknown unit "m," in duration "2m,2h"`,
		},
		{
			name:           "defaults returned",
			input:          []string{},
			flagDefault:    []time.Duration{time.Second, time.Minute},
			expectedValues: []time.Duration{time.Second, time.Minute},
		},
		{
			name:           "defaults overwritten",
			input:          []string{"1m", "1s"},
			flagDefault:    []time.Duration{time.Second, time.Minute},
			expectedValues: []time.Duration{time.Minute, time.Second},
		},
		{
			name:  "replace values",
			input: []string{"1m", "1h"},
			visitor: func(f *zflag.Flag) {
				if val, ok := f.Value.(zflag.SliceValue); ok {
					_ = val.Replace([]string{"1s"})
				}
			},
			expectedValues: []time.Duration{time.Second},
		},
		{
			name:           "trims input",
			input:          []string{" 1ns", "2ms  ", "  3m    ", "    4h"},
			expectedValues: []time.Duration{1 * time.Nanosecond, 2 * time.Millisecond, 3 * time.Minute, 4 * time.Hour},
		},
		{
			name:           "valid values",
			input:          []string{"1ns", "2ms", "3m", "4h"},
			expectedValues: []time.Duration{1 * time.Nanosecond, 2 * time.Millisecond, 3 * time.Minute, 4 * time.Hour},
		},
	}

	t.Parallel()
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			var ds []time.Duration
			f := zflag.NewFlagSet("test", zflag.ContinueOnError)
			f.SetOutput(ioutil.Discard)
			f.DurationSliceVar(&ds, "ds", test.flagDefault, "usage")
			err := f.Parse(repeatFlag("--ds", test.input...))
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

			if !reflect.DeepEqual(test.expectedValues, ds) {
				t.Fatalf("expected %v with type %T but got %v with type %T ", test.expectedValues, test.expectedValues, ds, ds)
			}

			durSlice, err := f.GetDurationSlice("ds")
			if err != nil {
				t.Fatal("got an error from GetDurationSlice():", err)
			}
			if !reflect.DeepEqual(test.expectedValues, durSlice) {
				t.Fatalf("expected %v with type %T but got %v with type %T ", test.expectedValues, test.expectedValues, durSlice, durSlice)
			}

			durSliceGet, err := f.Get("ds")
			if err != nil {
				t.Fatal("got an error from Get():", err)
			}
			if !reflect.DeepEqual(durSliceGet, durSlice) {
				t.Fatalf("expected %v with type %T but got %v with type %T ", test.expectedValues, test.expectedValues, durSliceGet, durSliceGet)
			}
		})
	}
}
