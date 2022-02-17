// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/gowarden/zflag"
)

const expectedOutput = `      --long-form    Some description
      --long-form2   Some description
                       with multiline
  -s, --long-name    Some description
  -t, --long-name2   Some description with
                       multiline
`

func setUpZFlagSet(buf io.Writer) *zflag.FlagSet {
	f := zflag.NewFlagSet("test", zflag.ExitOnError)
	f.Bool("long-form", false, "Some description")
	f.Bool("long-form2", false, "Some description\n  with multiline")
	f.Bool("long-name", false, "Some description", zflag.OptShorthand('s'))
	f.Bool("long-name2", false, "Some description with\n  multiline", zflag.OptShorthand('t'))
	f.SetOutput(buf)
	return f
}

func TestPrintUsage(t *testing.T) {
	buf := bytes.Buffer{}
	f := setUpZFlagSet(&buf)
	f.PrintDefaults()
	res := buf.String()
	if res != expectedOutput {
		t.Errorf("Expected \n%s \nActual \n%s", expectedOutput, res)
	}
}

func setUpZFlagSet2(buf io.Writer) *zflag.FlagSet {
	f := zflag.NewFlagSet("test", zflag.ExitOnError)
	f.Bool("long-form", false, "Some description")
	f.Bool("long-form2", false, "Some description\n  with multiline")
	f.Bool("long-name", false, "Some description", zflag.OptShorthand('s'))
	f.Bool("long-name2", false, "Some description with\n  multiline", zflag.OptShorthand('t'))
	f.String("some-very-long-arg", "test", "Some very long description having break the limit", zflag.OptShorthand('l'))
	f.String("other-very-long-arg", "long-default-value", "Some very long description having break the limit", zflag.OptShorthand('o'))
	f.String("some-very-long-arg2", "very long default value", "Some very long description\nwith line break\nmultiple")
	f.SetOutput(buf)
	return f
}

const expectedOutput2 = `      --long-form                    Some description
      --long-form2                   Some description
                                       with multiline
  -s, --long-name                    Some description
  -t, --long-name2                   Some description with
                                       multiline
  -o, --other-very-long-arg string   Some very long description having
                                     break the limit (default
                                     "long-default-value")
  -l, --some-very-long-arg string    Some very long description having
                                     break the limit (default "test")
      --some-very-long-arg2 string   Some very long description
                                     with line break
                                     multiple (default "very long default
                                     value")
`

func TestPrintUsage_2(t *testing.T) {
	buf := bytes.Buffer{}
	f := setUpZFlagSet2(&buf)
	res := f.FlagUsagesWrapped(80)
	if res != expectedOutput2 {
		t.Errorf("Expected \n%q \nActual \n%q", expectedOutput2, res)
	}
}
