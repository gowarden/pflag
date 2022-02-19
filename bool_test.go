// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag_test

import (
	"reflect"
	"testing"

	"github.com/gowarden/zflag"
)

func TestBoolValueImplementsGetter(t *testing.T) {
	f := zflag.NewFlagSet("test", zflag.ContinueOnError)
	f.Bool("bool", false, "bool")
	v := f.Lookup("bool").Value

	if _, ok := v.(zflag.Getter); !ok {
		t.Fatalf("%T should implement the Getter interface", v)
	}
}

func TestBool(t *testing.T) {
	tests := []struct {
		name          string
		flagDefault   bool
		input         []string
		expectedErr   string
		expectedValue bool
	}{
		{
			name:          "no value passed",
			input:         []string{},
			flagDefault:   false,
			expectedErr:   "",
			expectedValue: false,
		},
		{
			name:        "empty value passed",
			input:       []string{""},
			flagDefault: false,
			expectedErr: `invalid argument "" for "--bs" flag: strconv.ParseBool: parsing "": invalid syntax`,
		},
		{
			name:        "invalid bool",
			input:       []string{"blabla"},
			flagDefault: false,
			expectedErr: `invalid argument "blabla" for "--bs" flag: strconv.ParseBool: parsing "blabla": invalid syntax`,
		},
		{
			name:        "no csv",
			input:       []string{"true,false"},
			flagDefault: false,
			expectedErr: `invalid argument "true,false" for "--bs" flag: strconv.ParseBool: parsing "true,false": invalid syntax`,
		},
		{
			name:          "repeated value",
			input:         []string{"true", "false"},
			flagDefault:   true,
			expectedValue: false,
		},
		{
			name:          "with default values",
			input:         []string{"false"},
			flagDefault:   true,
			expectedValue: false,
		},
		{
			name:          "trims input true",
			input:         []string{" true "},
			expectedValue: true,
		},
		{
			name:          "trims input false",
			input:         []string{" false "},
			expectedValue: false,
		},
		{
			name:          "test all valid bools",
			input:         []string{"true"},
			expectedValue: true,
		},
		{
			name:          "test all valid bools",
			input:         []string{"false"},
			expectedValue: false,
		},
		{
			name:          "test all valid bools",
			input:         []string{"1"},
			expectedValue: true,
		},
		{
			name:          "test all valid bools",
			input:         []string{"0"},
			expectedValue: false,
		},
		{
			name:          "test all valid bools",
			input:         []string{"t"},
			expectedValue: true,
		},
		{
			name:          "test all valid bools",
			input:         []string{"f"},
			expectedValue: false,
		},
		{
			name:          "test all valid bools",
			input:         []string{"TRUE"},
			expectedValue: true,
		},
		{
			name:          "test all valid bools",
			input:         []string{"FALSE"},
			expectedValue: false,
		},
		{
			name:          "test all valid bools",
			input:         []string{"1"},
			expectedValue: true,
		},
		{
			name:          "test all valid bools",
			input:         []string{"0"},
			expectedValue: false,
		},
		{
			name:          "test all valid bools",
			input:         []string{"T"},
			expectedValue: true,
		},
		{
			name:          "test all valid bools",
			input:         []string{"F"},
			expectedValue: false,
		},
		{
			name:          "test all valid bools",
			input:         []string{"True"},
			expectedValue: true,
		},
		{
			name:          "test all valid bools",
			input:         []string{"False"},
			expectedValue: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var bs bool
			f := zflag.NewFlagSet("test", zflag.ContinueOnError)
			f.BoolVar(&bs, "bs", test.flagDefault, "usage")
			err := f.Parse(repeatFlag("--bs", test.input...))
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

			if !reflect.DeepEqual(test.expectedValue, bs) {
				t.Fatalf("expected %v with type %T but got %v with type %T ", test.expectedValue, test.expectedValue, bs, bs)
			}

			getBS, err := f.GetBool("bs")
			if err != nil {
				t.Fatal("got an error from GetBool():", err)
			}
			if !reflect.DeepEqual(test.expectedValue, getBS) {
				t.Fatalf("expected %v with type %T but got %v with type %T ", test.expectedValue, test.expectedValue, getBS, getBS)
			}

			getBSGet, err := f.Get("bs")
			if err != nil {
				t.Fatal("got an error from Get():", err)
			}
			if !reflect.DeepEqual(getBSGet, getBS) {
				t.Fatalf("expected %v with type %T but got %v with type %T ", test.expectedValue, test.expectedValue, getBS, getBS)
			}
		})
	}
}
