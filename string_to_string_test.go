// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag_test

import (
	"reflect"
	"testing"

	"github.com/gowarden/zflag"
)

func setUpS2SFlagSet(s2sp *map[string]string) *zflag.FlagSet {
	f := zflag.NewFlagSet("test", zflag.ContinueOnError)
	f.StringToStringVar(s2sp, "s2s", map[string]string{}, "Command separated ls2st!")
	return f
}

func setUpS2SFlagSetWithDefault(s2sp *map[string]string) *zflag.FlagSet {
	f := zflag.NewFlagSet("test", zflag.ContinueOnError)
	f.StringToStringVar(s2sp, "s2s", map[string]string{"da": "1", "db": "2", "de": "5,6"}, "Command separated ls2st!")
	return f
}

func createS2SFlag(vals map[string]string) []string {
	records := make([]string, 0, len(vals)>>1)
	for k, v := range vals {
		records = append(records, k+"="+v)
	}

	return records
}

func TestS2SValueImplementsGetter(t *testing.T) {
	f := zflag.NewFlagSet("test", zflag.ContinueOnError)
	f.StringToString("s2s", map[string]string{}, "Command separated ls2st!")
	v := f.Lookup("s2s").Value

	if _, ok := v.(zflag.Getter); !ok {
		t.Fatalf("%T should implement the Getter interface", v)
	}
}

func TestEmptyS2S(t *testing.T) {
	var s2s map[string]string
	f := setUpS2SFlagSet(&s2s)
	err := f.Parse([]string{})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}

	getS2S, err := f.GetStringToString("s2s")
	if err != nil {
		t.Fatal("got an error from GetStringToString():", err)
	}
	if len(getS2S) != 0 {
		t.Fatalf("got s2s %v with len=%d but expected length=0", getS2S, len(getS2S))
	}
	getS2S_2, err := f.Get("s2s")
	if err != nil {
		t.Fatal("got an error from Get():", err)
	}
	if !reflect.DeepEqual(getS2S_2, getS2S) {
		t.Fatalf("expected %v with type %T but got %v with type %T ", getS2S, getS2S, getS2S_2, getS2S_2)
	}
}

func TestS2S(t *testing.T) {
	var s2s map[string]string
	f := setUpS2SFlagSet(&s2s)

	vals := map[string]string{"a": "1", "b": "2", "d": "4", "c": "3", "e": "5,6"}
	err := f.Parse(repeatFlag("--s2s", createS2SFlag(vals)...))
	if err != nil {
		t.Fatal("expected no error; got", err)
	}
	for k, v := range s2s {
		if vals[k] != v {
			t.Fatalf("expected s2s[%s] to be %s but got: %s", k, vals[k], v)
		}
	}
	getS2S, err := f.GetStringToString("s2s")
	if err != nil {
		t.Fatalf("got error: %v", err)
	}
	for k, v := range getS2S {
		if vals[k] != v {
			t.Fatalf("expected s2s[%s] to be %s but got: %s from GetStringToString", k, vals[k], v)
		}
	}
	getS2S_2, err := f.Get("s2s")
	if err != nil {
		t.Fatal("got an error from Get():", err)
	}
	if !reflect.DeepEqual(getS2S_2, getS2S) {
		t.Fatalf("expected %v with type %T but got %v with type %T ", getS2S, getS2S, getS2S_2, getS2S_2)
	}
}

func TestS2SDefault(t *testing.T) {
	var s2s map[string]string
	f := setUpS2SFlagSetWithDefault(&s2s)

	vals := map[string]string{"da": "1", "db": "2", "de": "5,6"}

	err := f.Parse([]string{})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}
	for k, v := range s2s {
		if vals[k] != v {
			t.Fatalf("expected s2s[%s] to be %s but got: %s", k, vals[k], v)
		}
	}

	getS2S, err := f.GetStringToString("s2s")
	if err != nil {
		t.Fatal("got an error from GetStringToString():", err)
	}
	for k, v := range getS2S {
		if vals[k] != v {
			t.Fatalf("expected s2s[%s] to be %s from GetStringToString but got: %s", k, vals[k], v)
		}
	}
	getS2S_2, err := f.Get("s2s")
	if err != nil {
		t.Fatal("got an error from Get():", err)
	}
	if !reflect.DeepEqual(getS2S_2, getS2S) {
		t.Fatalf("expected %v with type %T but got %v with type %T ", getS2S, getS2S, getS2S_2, getS2S_2)
	}
}

func TestS2SWithDefault(t *testing.T) {
	var s2s map[string]string
	f := setUpS2SFlagSetWithDefault(&s2s)

	vals := map[string]string{"a": "1", "b": "2", "e": "5,6"}
	err := f.Parse(repeatFlag("--s2s", createS2SFlag(vals)...))
	if err != nil {
		t.Fatal("expected no error; got", err)
	}
	for k, v := range s2s {
		if vals[k] != v {
			t.Fatalf("expected s2s[%s] to be %s but got: %s", k, vals[k], v)
		}
	}

	getS2S, err := f.GetStringToString("s2s")
	if err != nil {
		t.Fatal("got an error from GetStringToString():", err)
	}
	for k, v := range getS2S {
		if vals[k] != v {
			t.Fatalf("expected s2s[%s] to be %s from GetStringToString but got: %s", k, vals[k], v)
		}
	}
	getS2S_2, err := f.Get("s2s")
	if err != nil {
		t.Fatal("got an error from Get():", err)
	}
	if !reflect.DeepEqual(getS2S_2, getS2S) {
		t.Fatalf("expected %v with type %T but got %v with type %T ", getS2S, getS2S, getS2S_2, getS2S_2)
	}
}

func TestS2SCalledTwice(t *testing.T) {
	var s2s map[string]string
	f := setUpS2SFlagSet(&s2s)

	in := []string{"a=1,b=2", "b=3", `e=5,6`}
	expected := map[string]string{"a": "1,b=2", "b": "3", "e": "5,6"}
	err := f.Parse(repeatFlag("--s2s", in...))
	if err != nil {
		t.Fatal("expected no error; got", err)
	}
	if len(s2s) != len(expected) {
		t.Fatalf("expected %d flags; got %d flags", len(expected), len(s2s))
	}
	for i, v := range s2s {
		if expected[i] != v {
			t.Fatalf("expected s2s[%s] to be %s but got: %s", i, expected[i], v)
		}
	}
}
