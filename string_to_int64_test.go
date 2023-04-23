// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag_test

import (
	"io/ioutil"
	"sort"
	"strings"
	"testing"

	"github.com/zulucmd/zflag/v2"
)

func TestStringToInt64(t *testing.T) {
	tests := []struct {
		name              string
		flagDefault       map[string]int64
		input             []string
		expectedErr       string
		expectedValues    map[string]int64
		expectedStrValues []string
		visitor           func(f *zflag.Flag)
	}{
		{
			name:              "no value passed",
			input:             []string{},
			flagDefault:       map[string]int64{},
			expectedErr:       "",
			expectedValues:    map[string]int64{},
			expectedStrValues: []string{},
		},
		{
			name:        "empty value passed",
			input:       []string{""},
			flagDefault: map[string]int64{},
			expectedErr: `invalid argument "" for "--s2i64" flag:  must be formatted as key=value`,
		},
		{
			name:        "invalid int64",
			input:       []string{"blabla"},
			flagDefault: map[string]int64{},
			expectedErr: `invalid argument "blabla" for "--s2i64" flag: blabla must be formatted as key=value`,
		},
		{
			name:        "no csv",
			input:       []string{"test=1,5"},
			flagDefault: map[string]int64{},
			expectedErr: `invalid argument "test=1,5" for "--s2i64" flag: strconv.ParseInt: parsing "1,5": invalid syntax`,
		},
		{
			name:        "single key value pair per arg",
			input:       []string{"test=1=1"},
			flagDefault: map[string]int64{},
			expectedErr: `invalid argument "test=1=1" for "--s2i64" flag: strconv.ParseInt: parsing "1=1": invalid syntax`,
		},
		{
			name:              "overrides multiple calls",
			input:             []string{"test=1", "test=5"},
			flagDefault:       map[string]int64{},
			expectedValues:    map[string]int64{"test": 5},
			expectedStrValues: []string{"test=5"},
		},
		{
			name:              "empty defaults",
			input:             []string{"test=1", "test2=5"},
			flagDefault:       map[string]int64{},
			expectedValues:    map[string]int64{"test": 1, "test2": 5},
			expectedStrValues: []string{"test=1", "test2=5"},
		},
		{
			name:              "overrides default values",
			input:             []string{"test=1", "test2=5"},
			flagDefault:       map[string]int64{"test2": 1, "test": 5},
			expectedValues:    map[string]int64{"test": 1, "test2": 5},
			expectedStrValues: []string{"test=1", "test2=5"},
		},
		{
			name:              "returns default values",
			input:             []string{},
			flagDefault:       map[string]int64{"test2": 1, "test": 5},
			expectedValues:    map[string]int64{"test2": 1, "test": 5},
			expectedStrValues: []string{"test2=1", "test=5"},
		},
		{
			name:              "trims input",
			input:             []string{"test=    1", "test2=5     ", "test3=     9     "},
			flagDefault:       map[string]int64{},
			expectedValues:    map[string]int64{"test": 1, "test2": 5, "test3": 9},
			expectedStrValues: []string{"test=1", "test2=5", "test3=9"},
		},
	}

	t.Parallel()
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			var s2i64 map[string]int64
			f := zflag.NewFlagSet("test", zflag.ContinueOnError)
			f.SetOutput(ioutil.Discard)
			f.StringToInt64Var(&s2i64, "s2i64", test.flagDefault, "usage")
			err := f.Parse(repeatFlag("--s2i64", test.input...))
			if test.expectedErr != "" {
				assertErr(t, err)
				assertEqualf(t, test.expectedErr, err.Error(), "expected error to equal %q, but was: %s", test.expectedErr, err)
				return
			}

			assertNoErr(t, err)

			if test.visitor != nil {
				f.VisitAll(test.visitor)
			}

			assertDeepEqual(t, test.expectedValues, s2i64)

			s2i64GetS2I64, err := f.GetStringToInt64("s2i64")
			assertNoErr(t, err)
			assertDeepEqual(t, test.expectedValues, s2i64GetS2I64)

			s2i64Get, err := f.Get("s2i64")
			assertNoErr(t, err)
			assertDeepEqual(t, test.expectedValues, s2i64Get)

			flag := f.Lookup("s2i64")
			strVal := flag.Value.String()
			if len(test.expectedStrValues) == 0 {
				assertEqual(t, "[]", strVal)
			} else {
				assertEqual(t, '[', rune(strVal[0]))
				assertEqual(t, ']', rune(strVal[len(strVal)-1]))

				strVals := strings.Split(strVal[1:len(strVal)-1], " ")
				sort.Strings(strVals)
				sort.Strings(test.expectedStrValues)
				assertDeepEqual(t, test.expectedStrValues, strVals)
			}

			defer assertNoPanic(t)()
			mustStringToInt64 := f.MustGetStringToInt64("s2i64")
			assertDeepEqual(t, test.expectedValues, mustStringToInt64)
		})
	}
}

func TestStringToInt64Errors(t *testing.T) {
	t.Parallel()

	var s string
	var s2i64 map[string]int64
	f := zflag.NewFlagSet("test", zflag.ContinueOnError)
	f.SetOutput(ioutil.Discard)
	f.StringVar(&s, "s", "", "usage")
	f.StringToInt64Var(&s2i64, "s2i64", map[string]int64{}, "usage")
	err := f.Parse([]string{})
	assertNoErr(t, err)

	_, err = f.GetStringToInt64("s")
	assertErr(t, err)

	defer assertPanic(t)()
	_ = f.MustGetStringToInt64("s")
}
