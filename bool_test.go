// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag_test

import (
	"bytes"
	"strconv"
	"testing"

	"github.com/gowarden/zflag"
)

// This value can be a boolean ("true", "false") or "maybe"
type triStateValue int

const (
	triStateFalse triStateValue = 0
	triStateTrue  triStateValue = 1
	triStateMaybe triStateValue = 2
)

const strTriStateMaybe = "maybe"

func (v *triStateValue) IsBoolFlag() bool {
	return true
}

func (v *triStateValue) Get() interface{} {
	return triStateValue(*v)
}

func (v *triStateValue) Set(s string) error {
	if s == strTriStateMaybe {
		*v = triStateMaybe
		return nil
	}
	boolVal, err := strconv.ParseBool(s)
	if boolVal {
		*v = triStateTrue
	} else {
		*v = triStateFalse
	}
	return err
}

func (v *triStateValue) String() string {
	if *v == triStateMaybe {
		return strTriStateMaybe
	}
	return strconv.FormatBool(*v == triStateTrue)
}

// The type of the flag as required by the zflag.Value interface
func (v *triStateValue) Type() string {
	return "version"
}

func setUpFlagSet(tristate *triStateValue) *zflag.FlagSet {
	f := zflag.NewFlagSet("test", zflag.ContinueOnError)
	*tristate = triStateFalse
	f.Var(tristate, "tristate", "tristate value (true, maybe or false)", zflag.OptShorthand('t'), zflag.OptNoOptDefVal("true"))
	return f
}

func TestBoolValueImplementsGetter(t *testing.T) {
	f := zflag.NewFlagSet("test", zflag.ContinueOnError)
	f.Bool("bool", false, "bool")
	v := f.Lookup("bool").Value

	if _, ok := v.(zflag.Getter); !ok {
		t.Fatalf("%T should implement the Getter interface", v)
	}
}

func TestExplicitTrue(t *testing.T) {
	var tristate triStateValue
	f := setUpFlagSet(&tristate)
	err := f.Parse([]string{"--tristate=true"})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}
	if tristate != triStateTrue {
		t.Fatal("expected", triStateTrue, "(triStateTrue) but got", tristate, "instead")
	}
}

func TestImplicitTrue(t *testing.T) {
	var tristate triStateValue
	f := setUpFlagSet(&tristate)
	err := f.Parse([]string{"--tristate"})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}
	if tristate != triStateTrue {
		t.Fatal("expected", triStateTrue, "(triStateTrue) but got", tristate, "instead")
	}
}

func TestShortFlag(t *testing.T) {
	var tristate triStateValue
	f := setUpFlagSet(&tristate)
	err := f.Parse([]string{"-t"})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}
	if tristate != triStateTrue {
		t.Fatal("expected", triStateTrue, "(triStateTrue) but got", tristate, "instead")
	}
}

func TestShortFlagExtraArgument(t *testing.T) {
	var tristate triStateValue
	f := setUpFlagSet(&tristate)
	// The"maybe"turns into an arg, since short boolean options will only do true/false
	err := f.Parse([]string{"-t", "maybe"})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}
	if tristate != triStateTrue {
		t.Fatal("expected", triStateTrue, "(triStateTrue) but got", tristate, "instead")
	}
	args := f.Args()
	if len(args) != 1 || args[0] != "maybe" {
		t.Fatal("expected an extra 'maybe' argument to stick around")
	}
}

func TestExplicitMaybe(t *testing.T) {
	var tristate triStateValue
	f := setUpFlagSet(&tristate)
	err := f.Parse([]string{"--tristate=maybe"})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}
	if tristate != triStateMaybe {
		t.Fatal("expected", triStateMaybe, "(triStateMaybe) but got", tristate, "instead")
	}
}

func TestExplicitFalse(t *testing.T) {
	var tristate triStateValue
	f := setUpFlagSet(&tristate)
	err := f.Parse([]string{"--tristate=false"})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}
	if tristate != triStateFalse {
		t.Fatal("expected", triStateFalse, "(triStateFalse) but got", tristate, "instead")
	}
}

func TestImplicitFalse(t *testing.T) {
	var tristate triStateValue
	f := setUpFlagSet(&tristate)
	err := f.Parse([]string{})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}
	if tristate != triStateFalse {
		t.Fatal("expected", triStateFalse, "(triStateFalse) but got", tristate, "instead")
	}
}

func TestInvalidValue(t *testing.T) {
	var tristate triStateValue
	f := setUpFlagSet(&tristate)
	var buf bytes.Buffer
	f.SetOutput(&buf)
	err := f.Parse([]string{"--tristate=invalid"})
	if err == nil {
		t.Fatal("expected an error but did not get any, tristate has value", tristate)
	}
}

func TestBoolP(t *testing.T) {
	b := zflag.Bool("bool", false, "bool value in CommandLine", zflag.OptShorthand('b'))
	c := zflag.Bool("c", false, "other bool value", zflag.OptShorthand('c'))
	args := []string{"--bool"}
	if err := zflag.CommandLine.Parse(args); err != nil {
		t.Error("expected no error, got ", err)
	}
	if *b != true {
		t.Errorf("expected b=true got b=%v", *b)
	}
	if *c != false {
		t.Errorf("expect c=false got c=%v", *c)
	}
}
