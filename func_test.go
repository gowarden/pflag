// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag

import (
	"fmt"
	"strings"
	"testing"
)

// Copyright 2009 The Go Authors. All rights reserved.
func TestUserDefinedFunc(t *testing.T) {
	var flags FlagSet
	flags.Init("test", ContinueOnError)
	var ss []string
	flags.Func("v", "usage", func(s string) error {
		ss = append(ss, s)
		return nil
	})
	if err := flags.Parse([]string{"-v", "1", "-v", "2", "-v=3"}); err != nil {
		t.Fatal(err)
	}
	if len(ss) != 3 {
		t.Fatal("expected 3 args; got ", len(ss))
	}
	expect := "[1 2 3]"
	if got := fmt.Sprint(ss); got != expect {
		t.Errorf("expected value %q got %q", expect, got)
	}
	// test usage
	var buf strings.Builder
	flags.SetOutput(&buf)
	flags.Parse([]string{"-h"})
	if usage := buf.String(); !strings.Contains(usage, "usage") {
		t.Errorf("usage string not included: %q", usage)
	}
	// test Func error
	flags = *NewFlagSet("test", ContinueOnError)
	flags.Func("v", "usage", func(s string) error {
		return fmt.Errorf("test error")
	})
	// flag not set, so no error
	if err := flags.Parse(nil); err != nil {
		t.Error(err)
	}
	// flag set, expect error
	if err := flags.Parse([]string{"-v", "1"}); err == nil {
		t.Error("expected error; got none")
	} else if errMsg := err.Error(); !strings.Contains(errMsg, "test error") {
		t.Errorf(`error should contain "test error"; got %q`, errMsg)
	}
}
