// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag_test

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"

	"github.com/gowarden/zflag"
)

func setUpS2I64FlagSet(s2ip *map[string]int64) *zflag.FlagSet {
	f := zflag.NewFlagSet("test", zflag.ContinueOnError)
	f.StringToInt64Var(s2ip, "s2i", map[string]int64{}, "Command separated ls2it!")
	return f
}

func setUpS2I64FlagSetWithDefault(s2ip *map[string]int64) *zflag.FlagSet {
	f := zflag.NewFlagSet("test", zflag.ContinueOnError)
	f.StringToInt64Var(s2ip, "s2i", map[string]int64{"a": 1, "b": 2}, "Command separated ls2it!")
	return f
}

func createS2I64Flag(vals map[string]int64) []string {
	var r []string
	for k, v := range vals {
		r = append(r, fmt.Sprintf("%s=%s", k, strconv.FormatInt(v, 10)))
	}
	return r
}

func TestS2I64ValueImplementsGetter(t *testing.T) {
	f := zflag.NewFlagSet("test", zflag.ContinueOnError)
	f.StringToInt64("s2i", map[string]int64{}, "Command separated ls2it!")
	v := f.Lookup("s2i").Value

	if _, ok := v.(zflag.Getter); !ok {
		t.Fatalf("%T should implement the Getter interface", v)
	}
}

func TestEmptyS2I64(t *testing.T) {
	var s2i map[string]int64
	f := setUpS2I64FlagSet(&s2i)
	err := f.Parse([]string{})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}

	getS2I, err := f.GetStringToInt64("s2i")
	if err != nil {
		t.Fatal("got an error from GetStringToInt64():", err)
	}
	if len(getS2I) != 0 {
		t.Fatalf("got s2i %v with len=%d but expected length=0", getS2I, len(getS2I))
	}
	getS2I_2, err := f.Get("s2i")
	if err != nil {
		t.Fatal("got an error from Get():", err)
	}
	if !reflect.DeepEqual(getS2I_2, getS2I) {
		t.Fatalf("expected %v with type %T but got %v with type %T ", getS2I_2, getS2I, getS2I_2, getS2I_2)
	}
}

func TestS2I64(t *testing.T) {
	var s2i map[string]int64
	f := setUpS2I64FlagSet(&s2i)

	vals := map[string]int64{"a": 1, "b": 2, "d": 4, "c": 3}
	err := f.Parse(repeatFlag("--s2i", createS2I64Flag(vals)...))
	if err != nil {
		t.Fatal("expected no error; got", err)
	}
	for k, v := range s2i {
		if vals[k] != v {
			t.Fatalf("expected s2i[%s] to be %d but got: %d", k, vals[k], v)
		}
	}
	getS2I, err := f.GetStringToInt64("s2i")
	if err != nil {
		t.Fatalf("got error: %v", err)
	}
	for k, v := range getS2I {
		if vals[k] != v {
			t.Fatalf("expected s2i[%s] to be %d but got: %d from GetStringToInt64", k, vals[k], v)
		}
	}
	getS2I_2, err := f.Get("s2i")
	if err != nil {
		t.Fatal("got an error from Get():", err)
	}
	if !reflect.DeepEqual(getS2I_2, getS2I) {
		t.Fatalf("expected %v with type %T but got %v with type %T ", getS2I, getS2I, getS2I_2, getS2I_2)
	}
}

func TestS2I64Default(t *testing.T) {
	var s2i map[string]int64
	f := setUpS2I64FlagSetWithDefault(&s2i)

	vals := map[string]int64{"a": 1, "b": 2}

	err := f.Parse([]string{})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}
	for k, v := range s2i {
		if vals[k] != v {
			t.Fatalf("expected s2i[%s] to be %d but got: %d", k, vals[k], v)
		}
	}

	getS2I, err := f.GetStringToInt64("s2i")
	if err != nil {
		t.Fatal("got an error from GetStringToInt64():", err)
	}
	for k, v := range getS2I {
		if vals[k] != v {
			t.Fatalf("expected s2i[%s] to be %d from GetStringToInt64 but got: %d", k, vals[k], v)
		}
	}
	getS2I_2, err := f.Get("s2i")
	if err != nil {
		t.Fatal("got an error from Get():", err)
	}
	if !reflect.DeepEqual(getS2I_2, getS2I) {
		t.Fatalf("expected %v with type %T but got %v with type %T ", getS2I, getS2I, getS2I_2, getS2I_2)
	}
}

func TestS2I64WithDefault(t *testing.T) {
	var s2i map[string]int64
	f := setUpS2I64FlagSetWithDefault(&s2i)

	vals := map[string]int64{"a": 1, "b": 2}
	err := f.Parse(repeatFlag("--s2i", createS2I64Flag(vals)...))
	if err != nil {
		t.Fatal("expected no error; got", err)
	}
	for k, v := range s2i {
		if vals[k] != v {
			t.Fatalf("expected s2i[%s] to be %d but got: %d", k, vals[k], v)
		}
	}

	getS2I, err := f.GetStringToInt64("s2i")
	if err != nil {
		t.Fatal("got an error from GetStringToInt64():", err)
	}
	for k, v := range getS2I {
		if vals[k] != v {
			t.Fatalf("expected s2i[%s] to be %d from GetStringToInt64 but got: %d", k, vals[k], v)
		}
	}
	getS2I_2, err := f.Get("s2i")
	if err != nil {
		t.Fatal("got an error from Get():", err)
	}
	if !reflect.DeepEqual(getS2I_2, getS2I) {
		t.Fatalf("expected %v with type %T but got %v with type %T ", getS2I, getS2I, getS2I_2, getS2I_2)
	}
}

func TestS2I64CalledTwice(t *testing.T) {
	var s2i map[string]int64
	f := setUpS2I64FlagSet(&s2i)

	in := []string{"a=1", "b=2", "b=3"}
	expected := map[string]int64{"a": 1, "b": 3}
	err := f.Parse(repeatFlag("--s2i", in...))
	if err != nil {
		t.Fatal("expected no error; got", err)
	}
	for i, v := range s2i {
		if expected[i] != v {
			t.Fatalf("expected s2i[%s] to be %d but got: %d", i, expected[i], v)
		}
	}
}
