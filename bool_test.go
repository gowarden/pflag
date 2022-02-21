// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag_test

import (
	"io/ioutil"
	"testing"

	"github.com/gowarden/zflag"
)

func TestBool(t *testing.T) {
	tests := []struct {
		name          string
		flagDefault   bool
		input         []string
		expectedErr   string
		expectedValue bool
		extraOpts     []zflag.Opt
	}{
		{
			name:          "no value passed",
			input:         []string{},
			flagDefault:   false,
			expectedErr:   "",
			expectedValue: false,
		},
		{
			name:          "empty value passed",
			input:         repeatFlag("--bs", ""),
			flagDefault:   false,
			expectedValue: true,
		},
		{
			name:        "invalid bool",
			input:       repeatFlag("--bs", "blabla"),
			flagDefault: false,
			expectedErr: `invalid argument "blabla" for "--bs" flag: strconv.ParseBool: parsing "blabla": invalid syntax`,
		},
		{
			name:        "no csv",
			input:       repeatFlag("--bs", "true,false"),
			flagDefault: false,
			expectedErr: `invalid argument "true,false" for "--bs" flag: strconv.ParseBool: parsing "true,false": invalid syntax`,
		},
		{
			name:        "flag prefixed with `no-` is not found",
			input:       []string{"--no-bs"},
			flagDefault: true,
			expectedErr: "unknown flag: --no-bs",
			extraOpts:   []zflag.Opt{zflag.OptShorthand('b')},
		},
		{
			name:          "flag prefixed with `no-` is found when negative enabled",
			input:         []string{"--no-bs"},
			flagDefault:   true,
			expectedValue: false,
			extraOpts:     []zflag.Opt{zflag.OptShorthand('b'), zflag.OptAddNegative()},
		},
		{
			name:        "no- prefix does not accept a value",
			input:       []string{"--no-bs=true"},
			expectedErr: "flag cannot have a value: --no-bs=true",
			extraOpts:   []zflag.Opt{zflag.OptShorthand('b'), zflag.OptAddNegative()},
		},
		{
			name:          "accepts separate value without no-",
			input:         []string{"--bs", "true"},
			flagDefault:   false,
			expectedValue: true,
			extraOpts:     []zflag.Opt{zflag.OptShorthand('b'), zflag.OptAddNegative()},
		},
		{
			name:          "accepts value without no-",
			input:         []string{"--bs=true"},
			flagDefault:   false,
			expectedValue: true,
			extraOpts:     []zflag.Opt{zflag.OptShorthand('b'), zflag.OptAddNegative()},
		},
		{
			name:          "short opt",
			input:         []string{"-b"},
			flagDefault:   false,
			expectedValue: true,
			extraOpts:     []zflag.Opt{zflag.OptShorthand('b'), zflag.OptAddNegative()},
		},
		{
			name:          "short opt with value false",
			input:         []string{"-b=0"},
			flagDefault:   true,
			expectedValue: false,
			extraOpts:     []zflag.Opt{zflag.OptShorthand('b'), zflag.OptAddNegative()},
		},
		{
			name:          "short opt with value",
			input:         []string{"-b=1"},
			flagDefault:   false,
			expectedValue: true,
			extraOpts:     []zflag.Opt{zflag.OptShorthand('b'), zflag.OptAddNegative()},
		},
		{
			name:          "short opt with value",
			input:         []string{"-b=true"},
			flagDefault:   false,
			expectedValue: true,
			extraOpts:     []zflag.Opt{zflag.OptShorthand('b'), zflag.OptAddNegative()},
		},
		{
			name:          "repeated value",
			input:         repeatFlag("--bs", "true", "false"),
			flagDefault:   true,
			expectedValue: false,
		},
		{
			name:          "with default values",
			input:         repeatFlag("--bs", "false"),
			flagDefault:   true,
			expectedValue: false,
		},
		{
			name:          "trims input true",
			input:         repeatFlag("--bs", " true "),
			expectedValue: true,
		},
		{
			name:          "trims input false",
			input:         repeatFlag("--bs", " false "),
			expectedValue: false,
		},
		{
			name:          "test all valid bools",
			input:         repeatFlag("--bs", "true"),
			expectedValue: true,
		},
		{
			name:          "test all valid bools",
			input:         repeatFlag("--bs", "false"),
			expectedValue: false,
		},
		{
			name:          "test all valid bools",
			input:         repeatFlag("--bs", "1"),
			expectedValue: true,
		},
		{
			name:          "test all valid bools",
			input:         repeatFlag("--bs", "0"),
			expectedValue: false,
		},
		{
			name:          "test all valid bools",
			input:         repeatFlag("--bs", "t"),
			expectedValue: true,
		},
		{
			name:          "test all valid bools",
			input:         repeatFlag("--bs", "f"),
			expectedValue: false,
		},
		{
			name:          "test all valid bools",
			input:         repeatFlag("--bs", "true"),
			expectedValue: true,
		},
		{
			name:          "test all valid bools",
			input:         repeatFlag("--bs", "false"),
			expectedValue: false,
		},
		{
			name:          "test all valid bools",
			input:         repeatFlag("--bs", "1"),
			expectedValue: true,
		},
		{
			name:          "test all valid bools",
			input:         repeatFlag("--bs", "0"),
			expectedValue: false,
		},
		{
			name:          "test all valid bools",
			input:         repeatFlag("--bs", "t"),
			expectedValue: true,
		},
		{
			name:          "test all valid bools",
			input:         repeatFlag("--bs", "f"),
			expectedValue: false,
		},
		{
			name:          "test all valid bools",
			input:         repeatFlag("--bs", "true"),
			expectedValue: true,
		},
		{
			name:          "test all valid bools",
			input:         repeatFlag("--bs", "false"),
			expectedValue: false,
		},
	}

	t.Parallel()
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			var bs bool
			f := zflag.NewFlagSet("test", zflag.ContinueOnError)
			f.SetOutput(ioutil.Discard)
			f.BoolVar(&bs, "bs", test.flagDefault, "usage", test.extraOpts...)
			err := f.Parse(test.input)
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

			if bs != test.expectedValue {
				t.Fatalf("expected %v with type %T but got %v with type %T ", test.expectedValue, test.expectedValue, bs, bs)
			}

			getBS, err := f.GetBool("bs")
			if err != nil {
				t.Fatal("got an error from GetBool():", err)
			}
			if getBS != test.expectedValue {
				t.Fatalf("expected %v with type %T but got %v with type %T ", test.expectedValue, test.expectedValue, getBS, getBS)
			}

			getBSGet, err := f.Get("bs")
			if err != nil {
				t.Fatal("got an error from Get():", err)
			}
			if getBSGet != test.expectedValue {
				t.Fatalf("expected %v with type %T but got %v with type %T ", test.expectedValue, test.expectedValue, getBSGet, getBSGet)
			}
		})
	}
}
