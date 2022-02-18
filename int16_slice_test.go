// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/gowarden/zflag"
)

func setUpI16SFlagSet(isp *[]int16) *zflag.FlagSet {
	f := zflag.NewFlagSet("test", zflag.ContinueOnError)
	f.Int16SliceVar(isp, "is", []int16{}, "Command separated list!")
	return f
}

func setUpI16SFlagSetWithDefault(isp *[]int16) *zflag.FlagSet {
	f := zflag.NewFlagSet("test", zflag.ContinueOnError)
	f.Int16SliceVar(isp, "is", []int16{0, 1}, "Command separated list!")
	return f
}

func TestI16SliceValueImplementsGetter(t *testing.T) {
	f := zflag.NewFlagSet("test", zflag.ContinueOnError)
	f.Int16Slice("is", []int16{0, 1}, "Command separated list!")
	v := f.Lookup("is").Value

	if _, ok := v.(zflag.Getter); !ok {
		t.Fatalf("%T should implement the Getter interface", v)
	}
}

func TestEmptyI16S(t *testing.T) {
	var is []int16
	f := setUpI16SFlagSet(&is)
	err := f.Parse([]string{})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}

	getI16S, err := f.GetInt16Slice("is")
	if err != nil {
		t.Fatal("got an error from GetInt16Slice():", err)
	}
	if len(getI16S) != 0 {
		t.Fatalf("got is %v with len=%d but expected length=0", getI16S, len(getI16S))
	}
	getI16S2, err := f.Get("is")
	if err != nil {
		t.Fatal("got an error from Get():", err)
	}
	if len(getI16S2.([]int16)) != 0 {
		t.Fatalf("got bs %v with len=%d but expected length=0", getI16S2.([]int16), len(getI16S2.([]int16)))
	}
}

func TestI16S(t *testing.T) {
	var is []int16
	f := setUpI16SFlagSet(&is)

	vals := []string{"1", "2", "4", "3"}
	err := f.Parse(repeatFlag("--is", vals...))
	if err != nil {
		t.Fatal("expected no error; got", err)
	}
	for i, v := range is {
		d64, err := strconv.ParseInt(vals[i], 0, 16)
		if err != nil {
			t.Fatalf("got error: %v", err)
		}
		d := int16(d64)
		if d != v {
			t.Fatalf("expected is[%d] to be %s but got: %d", i, vals[i], v)
		}
	}
	getI16S, err := f.GetInt16Slice("is")
	if err != nil {
		t.Fatalf("got error: %v", err)
	}
	for i, v := range getI16S {
		d64, err := strconv.ParseInt(vals[i], 0, 16)
		if err != nil {
			t.Fatalf("got error: %v", err)
		}
		d := int16(d64)
		if d != v {
			t.Fatalf("expected is[%d] to be %s but got: %d from GetInt16Slice", i, vals[i], v)
		}
	}
}

func TestI16SDefault(t *testing.T) {
	var is []int16
	f := setUpI16SFlagSetWithDefault(&is)

	vals := []string{"0", "1"}

	err := f.Parse([]string{})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}
	for i, v := range is {
		d64, err := strconv.ParseInt(vals[i], 0, 16)
		if err != nil {
			t.Fatalf("got error: %v", err)
		}
		d := int16(d64)
		if d != v {
			t.Fatalf("expected is[%d] to be %d but got: %d", i, d, v)
		}
	}

	getI16S, err := f.GetInt16Slice("is")
	if err != nil {
		t.Fatal("got an error from GetInt16Slice():", err)
	}
	for i, v := range getI16S {
		d64, err := strconv.ParseInt(vals[i], 0, 16)
		if err != nil {
			t.Fatal("got an error from GetInt16Slice():", err)
		}
		d := int16(d64)
		if d != v {
			t.Fatalf("expected is[%d] to be %d from GetInt16Slice but got: %d", i, d, v)
		}
	}
}

func TestI16SWithDefault(t *testing.T) {
	var is []int16
	f := setUpI16SFlagSetWithDefault(&is)

	vals := []string{"1", "2"}
	err := f.Parse(repeatFlag("--is", vals...))
	if err != nil {
		t.Fatal("expected no error; got", err)
	}
	for i, v := range is {
		d64, err := strconv.ParseInt(vals[i], 0, 16)
		if err != nil {
			t.Fatalf("got error: %v", err)
		}
		d := int16(d64)
		if d != v {
			t.Fatalf("expected is[%d] to be %d but got: %d", i, d, v)
		}
	}

	getI16S, err := f.GetInt16Slice("is")
	if err != nil {
		t.Fatal("got an error from GetInt16Slice():", err)
	}
	for i, v := range getI16S {
		d64, err := strconv.ParseInt(vals[i], 0, 16)
		if err != nil {
			t.Fatalf("got error: %v", err)
		}
		d := int16(d64)
		if d != v {
			t.Fatalf("expected is[%d] to be %d from GetInt16Slice but got: %d", i, d, v)
		}
	}
}

func TestI16SAsSliceValue(t *testing.T) {
	var i16s []int16
	f := setUpI16SFlagSet(&i16s)

	in := []string{"1", "2"}
	argfmt := "--is=%s"
	arg1 := fmt.Sprintf(argfmt, in[0])
	arg2 := fmt.Sprintf(argfmt, in[1])
	err := f.Parse([]string{arg1, arg2})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}

	f.VisitAll(func(f *zflag.Flag) {
		if val, ok := f.Value.(zflag.SliceValue); ok {
			_ = val.Replace([]string{"3"})
		}
	})
	if len(i16s) != 1 || i16s[0] != 3 {
		t.Fatalf("Expected ss to be overwritten with '3.1', but got: %v", i16s)
	}
}
