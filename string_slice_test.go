// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag_test

import (
	"reflect"
	"testing"

	"github.com/gowarden/zflag"
)

func setUpSSFlagSet(ssp *[]string) *zflag.FlagSet {
	f := zflag.NewFlagSet("test", zflag.ContinueOnError)
	f.StringSliceVar(ssp, "ss", []string{}, "Command separated list!")
	return f
}

func setUpSSFlagSetWithDefault(ssp *[]string) *zflag.FlagSet {
	f := zflag.NewFlagSet("test", zflag.ContinueOnError)
	f.StringSliceVar(ssp, "ss", []string{"default", "values"}, "Command separated list!")
	return f
}

func TestSSValueImplementsGetter(t *testing.T) {
	f := zflag.NewFlagSet("test", zflag.ContinueOnError)
	f.StringSlice("ss", []string{"default", "values"}, "Command separated list!")
	v := f.Lookup("ss").Value

	if _, ok := v.(zflag.Getter); !ok {
		t.Fatalf("%T should implement the Getter interface", v)
	}
}

func TestEmptySS(t *testing.T) {
	var ss []string
	f := setUpSSFlagSet(&ss)
	err := f.Parse([]string{})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}

	getSS, err := f.GetStringSlice("ss")
	if err != nil {
		t.Fatal("got an error from GetStringSlice():", err)
	}
	if len(getSS) != 0 {
		t.Fatalf("got ss %v with len=%d but expected length=0", getSS, len(getSS))
	}
	getSS_2, err := f.Get("ss")
	if err != nil {
		t.Fatal("got an error from Get():", err)
	}
	if !reflect.DeepEqual(getSS_2, getSS) {
		t.Fatalf("expected %v with type %T but got %v with type %T ", getSS, getSS, getSS_2, getSS_2)
	}
}

func TestEmptySSValue(t *testing.T) {
	var ss []string
	f := setUpSSFlagSet(&ss)
	err := f.Parse([]string{"--ss="})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}

	getSS, err := f.GetStringSlice("ss")
	if err != nil {
		t.Fatal("got an error from GetStringSlice():", err)
	}
	if len(getSS) != 1 {
		t.Fatalf("got ss %v with len=%d but expected length=1", getSS, len(getSS))
	}
	getSS_2, err := f.Get("ss")
	if err != nil {
		t.Fatal("got an error from Get():", err)
	}
	if !reflect.DeepEqual(getSS_2, getSS) {
		t.Fatalf("expected %v with type %T but got %v with type %T ", getSS, getSS, getSS_2, getSS_2)
	}
}

func TestSS(t *testing.T) {
	var ss []string
	f := setUpSSFlagSet(&ss)

	vals := []string{"one", "two", "4", "3"}
	err := f.Parse(repeatFlag("--ss", vals...))
	if err != nil {
		t.Fatal("expected no error; got", err)
	}
	for i, v := range ss {
		if vals[i] != v {
			t.Fatalf("expected ss[%d] to be %s but got: %s", i, vals[i], v)
		}
	}

	getSS, err := f.GetStringSlice("ss")
	if err != nil {
		t.Fatal("got an error from GetStringSlice():", err)
	}
	for i, v := range getSS {
		if vals[i] != v {
			t.Fatalf("expected ss[%d] to be %s from GetStringSlice but got: %s", i, vals[i], v)
		}
	}
	getSS_2, err := f.Get("ss")
	if err != nil {
		t.Fatal("got an error from Get():", err)
	}
	if !reflect.DeepEqual(getSS_2, getSS) {
		t.Fatalf("expected %v with type %T but got %v with type %T ", getSS, getSS, getSS_2, getSS_2)
	}
}

func TestSSDefault(t *testing.T) {
	var ss []string
	f := setUpSSFlagSetWithDefault(&ss)

	vals := []string{"default", "values"}

	err := f.Parse([]string{})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}
	for i, v := range ss {
		if vals[i] != v {
			t.Fatalf("expected ss[%d] to be %s but got: %s", i, vals[i], v)
		}
	}

	getSS, err := f.GetStringSlice("ss")
	if err != nil {
		t.Fatal("got an error from GetStringSlice():", err)
	}
	for i, v := range getSS {
		if vals[i] != v {
			t.Fatalf("expected ss[%d] to be %s from GetStringSlice but got: %s", i, vals[i], v)
		}
	}
	getSS_2, err := f.Get("ss")
	if err != nil {
		t.Fatal("got an error from Get():", err)
	}
	if !reflect.DeepEqual(getSS_2, getSS) {
		t.Fatalf("expected %v with type %T but got %v with type %T ", getSS, getSS, getSS_2, getSS_2)
	}
}

func TestSSWithDefault(t *testing.T) {
	var ss []string
	f := setUpSSFlagSetWithDefault(&ss)

	vals := []string{"one", "two", "4", "3"}
	err := f.Parse(repeatFlag("--ss", vals...))
	if err != nil {
		t.Fatal("expected no error; got", err)
	}
	for i, v := range ss {
		if vals[i] != v {
			t.Fatalf("expected ss[%d] to be %s but got: %s", i, vals[i], v)
		}
	}

	getSS, err := f.GetStringSlice("ss")
	if err != nil {
		t.Fatal("got an error from GetStringSlice():", err)
	}
	for i, v := range getSS {
		if vals[i] != v {
			t.Fatalf("expected ss[%d] to be %s from GetStringSlice but got: %s", i, vals[i], v)
		}
	}
	getSS_2, err := f.Get("ss")
	if err != nil {
		t.Fatal("got an error from Get():", err)
	}
	if !reflect.DeepEqual(getSS_2, getSS) {
		t.Fatalf("expected %v with type %T but got %v with type %T ", getSS, getSS, getSS_2, getSS_2)
	}
}

func TestSSCalledTwice(t *testing.T) {
	var ss []string
	f := setUpSSFlagSet(&ss)

	in := []string{"one", "two", "three"}
	expected := []string{"one", "two", "three"}
	err := f.Parse(repeatFlag("--ss", in...))
	if err != nil {
		t.Fatal("expected no error; got", err)
	}

	if len(expected) != len(ss) {
		t.Fatalf("expected number of ss to be %d but got: %d", len(expected), len(ss))
	}
	for i, v := range ss {
		if expected[i] != v {
			t.Fatalf("expected ss[%d] to be %s but got: %s", i, expected[i], v)
		}
	}

	values, err := f.GetStringSlice("ss")
	if err != nil {
		t.Fatal("expected no error; got", err)
	}

	if len(expected) != len(values) {
		t.Fatalf("expected number of values to be %d but got: %d", len(expected), len(ss))
	}
	for i, v := range values {
		if expected[i] != v {
			t.Fatalf("expected got ss[%d] to be %s but got: %s", i, expected[i], v)
		}
	}
	values_2, err := f.Get("ss")
	if err != nil {
		t.Fatal("got an error from Get():", err)
	}
	if !reflect.DeepEqual(values_2, values) {
		t.Fatalf("expected %v with type %T but got %v with type %T ", values, values, values_2, values_2)
	}
}

func TestSSNewLines(t *testing.T) {
	var got []string
	expected := []string{"foo\nbar\nbaz\n"}

	fs := zflag.NewFlagSet("test", zflag.ContinueOnError)
	fs.StringSliceVar(&got, "ss", []string{}, "")
	fs.Parse([]string{"--ss", expected[0]})
	if expected[0] != got[0] {
		t.Errorf("expected %q, got %q", expected[0], got[0])
	}
}
