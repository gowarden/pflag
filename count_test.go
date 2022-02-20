// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag_test

import (
	"io/ioutil"
	"testing"

	"github.com/gowarden/zflag"
)

func TestCountValueImplementsGetter(t *testing.T) {
	f := zflag.NewFlagSet("test", zflag.ContinueOnError)
	f.Count("verbose", "a counter")
	v := f.Lookup("verbose").Value

	if _, ok := v.(zflag.Getter); !ok {
		t.Fatalf("%T should implement the Getter interface", v)
	}
}

func TestCount(t *testing.T) {
	tests := []struct {
		name          string
		input         []string
		expectedErr   string
		expectedValue int
	}{
		{
			name:          "no flags",
			input:         []string{},
			expectedValue: 0,
		},
		{
			name:          "single count",
			input:         []string{"-v"},
			expectedValue: 1,
		},
		{
			name:          "multiple times",
			input:         []string{"-vvv"},
			expectedValue: 3,
		},
		{
			name:          "multiple times separated",
			input:         []string{"-v", "-v", "-v"},
			expectedValue: 3,
		},
		{
			name:          "multiple times interchanged and separated",
			input:         []string{"-v", "--verbose", "-v"},
			expectedValue: 3,
		},
		{
			name:          "multiple times with value",
			input:         []string{"-v=3", "-v"},
			expectedValue: 4,
		},
		{
			name:          "long opt with value",
			input:         []string{"--verbose=0"},
			expectedValue: 0,
		},
		{
			name:          "single with value",
			input:         []string{"-v=0"},
			expectedValue: 0,
		},
		{
			name:          "",
			input:         []string{"-v=a"},
			expectedErr:   `invalid argument "a" for "-v, --verbose" flag: strconv.ParseInt: parsing "a": invalid syntax`,
			expectedValue: 0,
		},
	}

	t.Parallel()
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			var verbose int
			f := zflag.NewFlagSet("test", zflag.ContinueOnError)
			f.SetOutput(ioutil.Discard)
			f.CountVar(&verbose, "verbose", "usage", zflag.OptShorthand('v'))
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

			if verbose != test.expectedValue {
				t.Fatalf("expected %v with type %T but got %v with type %T ", test.expectedValue, test.expectedValue, verbose, verbose)
			}

			getVerbose, err := f.GetCount("verbose")
			if err != nil {
				t.Fatal("got an error from GetCount():", err)
			}
			if getVerbose != test.expectedValue {
				t.Fatalf("expected %v with type %T but got %v with type %T ", test.expectedValue, test.expectedValue, getVerbose, getVerbose)
			}

			getVerboseGet, err := f.Get("verbose")
			if err != nil {
				t.Fatal("got an error from Get():", err)
			}
			if getVerboseGet != test.expectedValue {
				t.Fatalf("expected %v with type %T but got %v with type %T ", test.expectedValue, test.expectedValue, getVerboseGet, getVerboseGet)
			}
		})
	}
}
