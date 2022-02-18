// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag_test

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/gowarden/zflag"
)

func setUpISFlagSet(isp *[]int) *zflag.FlagSet {
	f := zflag.NewFlagSet("test", zflag.ContinueOnError)
	f.IntSliceVar(isp, "is", []int{}, "Command separated list!")
	return f
}

func setUpISFlagSetWithDefault(isp *[]int) *zflag.FlagSet {
	f := zflag.NewFlagSet("test", zflag.ContinueOnError)
	f.IntSliceVar(isp, "is", []int{0, 1}, "Command separated list!")
	return f
}

func TestISValueImplementsGetter(t *testing.T) {
	f := zflag.NewFlagSet("test", zflag.ContinueOnError)
	f.IntSlice("is", []int{0, 1}, "Command separated list!")
	v := f.Lookup("is").Value

	if _, ok := v.(zflag.Getter); !ok {
		t.Fatalf("%T should implement the Getter interface", v)
	}
}

func TestEmptyIS(t *testing.T) {
	var is []int
	f := setUpISFlagSet(&is)
	err := f.Parse([]string{})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}

	getIS, err := f.GetIntSlice("is")
	if err != nil {
		t.Fatal("got an error from GetIntSlice():", err)
	}
	if len(getIS) != 0 {
		t.Fatalf("got is %v with len=%d but expected length=0", getIS, len(getIS))
	}
	getIS_2, err := f.Get("is")
	if err != nil {
		t.Fatal("got an error from Get():", err)
	}

	if !reflect.DeepEqual(getIS_2, getIS) {
		t.Fatalf("expected %v with type %T but got %v with type %T ", getIS, getIS, getIS_2, getIS_2)
	}
}

func TestIS(t *testing.T) {
	var is []int
	f := setUpISFlagSet(&is)

	vals := []string{"1", "2", "4", "3"}
	err := f.Parse(repeatFlag("--is", vals...))
	if err != nil {
		t.Fatal("expected no error; got", err)
	}
	for i, v := range is {
		d, err := strconv.Atoi(vals[i])
		if err != nil {
			t.Fatalf("got error: %v", err)
		}
		if d != v {
			t.Fatalf("expected is[%d] to be %s but got: %d", i, vals[i], v)
		}
	}
	getIS, err := f.GetIntSlice("is")
	if err != nil {
		t.Fatalf("got error: %v", err)
	}
	for i, v := range getIS {
		d, err := strconv.Atoi(vals[i])
		if err != nil {
			t.Fatalf("got error: %v", err)
		}
		if d != v {
			t.Fatalf("expected is[%d] to be %s but got: %d from GetIntSlice", i, vals[i], v)
		}
	}
	getIS_2, err := f.Get("is")
	if err != nil {
		t.Fatal("got an error from Get():", err)
	}

	if !reflect.DeepEqual(getIS_2, getIS) {
		t.Fatalf("expected %v with type %T but got %v with type %T ", getIS, getIS, getIS_2, getIS_2)
	}
}

func TestISDefault(t *testing.T) {
	var is []int
	f := setUpISFlagSetWithDefault(&is)

	vals := []string{"0", "1"}

	err := f.Parse([]string{})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}
	for i, v := range is {
		d, err := strconv.Atoi(vals[i])
		if err != nil {
			t.Fatalf("got error: %v", err)
		}
		if d != v {
			t.Fatalf("expected is[%d] to be %d but got: %d", i, d, v)
		}
	}

	getIS, err := f.GetIntSlice("is")
	if err != nil {
		t.Fatal("got an error from GetIntSlice():", err)
	}
	for i, v := range getIS {
		d, err := strconv.Atoi(vals[i])
		if err != nil {
			t.Fatal("got an error from GetIntSlice():", err)
		}
		if d != v {
			t.Fatalf("expected is[%d] to be %d from GetIntSlice but got: %d", i, d, v)
		}
	}
}

func TestISWithDefault(t *testing.T) {
	var is []int
	f := setUpISFlagSetWithDefault(&is)

	vals := []string{"1", "2"}
	err := f.Parse(repeatFlag("--is", vals...))
	if err != nil {
		t.Fatal("expected no error; got", err)
	}
	for i, v := range is {
		d, err := strconv.Atoi(vals[i])
		if err != nil {
			t.Fatalf("got error: %v", err)
		}
		if d != v {
			t.Fatalf("expected is[%d] to be %d but got: %d", i, d, v)
		}
	}

	getIS, err := f.GetIntSlice("is")
	if err != nil {
		t.Fatal("got an error from GetIntSlice():", err)
	}
	for i, v := range getIS {
		d, err := strconv.Atoi(vals[i])
		if err != nil {
			t.Fatalf("got error: %v", err)
		}
		if d != v {
			t.Fatalf("expected is[%d] to be %d from GetIntSlice but got: %d", i, d, v)
		}
	}
}

func TestISCalledTwice(t *testing.T) {
	var is []int
	f := setUpISFlagSet(&is)

	in := []string{"1", "2", "3"}
	expected := []int{1, 2, 3}
	err := f.Parse(repeatFlag("--is", in...))
	if err != nil {
		t.Fatal("expected no error; got", err)
	}
	for i, v := range is {
		if expected[i] != v {
			t.Fatalf("expected is[%d] to be %d but got: %d", i, expected[i], v)
		}
	}
}
