// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/gowarden/zflag"
)

func setUpBSFlagSet(bsp *[]bool) *zflag.FlagSet {
	f := zflag.NewFlagSet("test", zflag.ContinueOnError)
	f.BoolSliceVar(bsp, "bs", []bool{}, "Command separated list!")
	return f
}

func setUpBSFlagSetWithDefault(bsp *[]bool) *zflag.FlagSet {
	f := zflag.NewFlagSet("test", zflag.ContinueOnError)
	f.BoolSliceVar(bsp, "bs", []bool{false, true}, "Command separated list!")
	return f
}

func TestBoolSliceValueImplementsGetter(t *testing.T) {
	f := zflag.NewFlagSet("test", zflag.ContinueOnError)
	f.BoolSlice("bs", []bool{false, true}, "Command separated list!")
	var v = f.Lookup("bs").Value

	if _, ok := v.(zflag.Getter); !ok {
		t.Fatalf("%T should implement the Getter interface", v)
	}
}

func TestEmptyBS(t *testing.T) {
	var bs []bool
	f := setUpBSFlagSet(&bs)
	err := f.Parse([]string{})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}

	getBS, err := f.GetBoolSlice("bs")
	if err != nil {
		t.Fatal("got an error from GetBoolSlice():", err)
	}
	if len(getBS) != 0 {
		t.Fatalf("got bs %v with len=%d but expected length=0", getBS, len(getBS))
	}
	getBS2, err := f.Get("bs")
	if err != nil {
		t.Fatal("got an error from Get():", err)
	}
	if len(getBS2.([]bool)) != 0 {
		t.Fatalf("got bs %v with len=%d but expected length=0", getBS2.([]bool), len(getBS2.([]bool)))
	}
}

func repeatFlag(flag string, values ...string) (res []string) {
	res = make([]string, 0, len(values))
	for _, val := range values {
		res = append(res, fmt.Sprintf("%s=%s", flag, val))
	}

	return
}

func TestBS(t *testing.T) {
	var bs []bool
	f := setUpBSFlagSet(&bs)

	vals := []string{"1", "F", "TRUE", "0"}
	err := f.Parse(repeatFlag("--bs", vals...))
	if err != nil {
		t.Fatal("expected no error; got", err)
	}
	for i, v := range bs {
		b, err := strconv.ParseBool(vals[i])
		if err != nil {
			t.Fatalf("got error: %v", err)
		}
		if b != v {
			t.Fatalf("expected is[%d] to be %s but got: %t", i, vals[i], v)
		}
	}
	getBS, err := f.GetBoolSlice("bs")
	if err != nil {
		t.Fatalf("got error: %v", err)
	}
	for i, v := range getBS {
		b, err := strconv.ParseBool(vals[i])
		if err != nil {
			t.Fatalf("got error: %v", err)
		}
		if b != v {
			t.Fatalf("expected bs[%d] to be %s but got: %t from GetBoolSlice", i, vals[i], v)
		}
	}
}

func TestBSDefault(t *testing.T) {
	var bs []bool
	f := setUpBSFlagSetWithDefault(&bs)

	vals := []string{"false", "T"}
	err := f.Parse([]string{})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}
	for i, v := range bs {
		b, err := strconv.ParseBool(vals[i])
		if err != nil {
			t.Fatalf("got error: %v", err)
		}
		if b != v {
			t.Fatalf("expected bs[%d] to be %t from GetBoolSlice but got: %t", i, b, v)
		}
	}

	getBS, err := f.GetBoolSlice("bs")
	if err != nil {
		t.Fatal("got an error from GetBoolSlice():", err)
	}
	for i, v := range getBS {
		b, err := strconv.ParseBool(vals[i])
		if err != nil {
			t.Fatal("got an error from GetBoolSlice():", err)
		}
		if b != v {
			t.Fatalf("expected bs[%d] to be %t from GetBoolSlice but got: %t", i, b, v)
		}
	}
}

func TestBSWithDefault(t *testing.T) {
	var bs []bool
	f := setUpBSFlagSetWithDefault(&bs)

	vals := []string{"FALSE", "1"}
	err := f.Parse(repeatFlag("--bs", vals...))
	if err != nil {
		t.Fatal("expected no error; got", err)
	}
	for i, v := range bs {
		b, err := strconv.ParseBool(vals[i])
		if err != nil {
			t.Fatalf("got error: %v", err)
		}
		if b != v {
			t.Fatalf("expected bs[%d] to be %t but got: %t", i, b, v)
		}
	}

	getBS, err := f.GetBoolSlice("bs")
	if err != nil {
		t.Fatal("got an error from GetBoolSlice():", err)
	}
	for i, v := range getBS {
		b, err := strconv.ParseBool(vals[i])
		if err != nil {
			t.Fatalf("got error: %v", err)
		}
		if b != v {
			t.Fatalf("expected bs[%d] to be %t from GetBoolSlice but got: %t", i, b, v)
		}
	}
}

func TestBSAsSliceValue(t *testing.T) {
	var bs []bool
	f := setUpBSFlagSet(&bs)

	in := []string{"true", "false"}
	err := f.Parse(repeatFlag("--bs", in...))
	if err != nil {
		t.Fatal("expected no error; got", err)
	}

	f.VisitAll(func(f *zflag.Flag) {
		if val, ok := f.Value.(zflag.SliceValue); ok {
			_ = val.Replace([]string{"false"})
		}
	})
	if len(bs) != 1 || bs[0] != false {
		t.Fatalf("Expected ss to be overwritten with 'false', but got: %v", bs)
	}
}

func TestBSBadSpacing(t *testing.T) {
	tests := []struct {
		Want    []bool
		FlagArg []string
	}{
		{
			Want:    []bool{true, false, true},
			FlagArg: []string{"1", "0", "true"},
		},
		{
			Want:    []bool{true, false},
			FlagArg: []string{"True", "F"},
		},
		{
			Want:    []bool{true, false},
			FlagArg: []string{"T", "0"},
		},
		{
			Want:    []bool{true, false},
			FlagArg: []string{"1", "0"},
		},
		{
			Want:    []bool{true, false, false},
			FlagArg: []string{"true", "false", "false"},
		},
		{
			Want:    []bool{true, false, false, true, false, true, false},
			FlagArg: []string{"true", "false", "false", "1", "0", "     T", " false "},
		},
		{
			Want:    []bool{false, false, true, false, true, false, true},
			FlagArg: []string{"0", " False", "  T", "false  ", " true", "F", "true"},
		},
	}

	for i, test := range tests {

		var bs []bool
		f := setUpBSFlagSet(&bs)

		if err := f.Parse(repeatFlag("--bs", test.FlagArg...)); err != nil {
			t.Fatalf("flag parsing failed with error: %s\nparsing:\t%#v\nwant:\t\t%#v",
				err, test.FlagArg, test.Want[i])
		}

		for j, b := range bs {
			if b != test.Want[j] {
				t.Fatalf("bad value parsed for test %d on bool %d:\nwant:\t%t\ngot:\t%t", i, j, test.Want[j], b)
			}
		}
	}
}
