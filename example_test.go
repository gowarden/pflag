// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag_test

import (
	"fmt"

	"github.com/gowarden/zflag"
)

func ExampleShorthandLookup() {
	name := "verbose"
	short := 'v'

	zflag.Bool(name, false, "verbose output", zflag.OptShorthand(short))

	// len(short) must be == 1
	flag := zflag.ShorthandLookup(short)

	fmt.Println(flag.Name)
}

func ExampleFlagSet_ShorthandLookup() {
	name := "verbose"
	short := 'v'

	fs := zflag.NewFlagSet("Example", zflag.ContinueOnError)
	fs.Bool(name, false, "verbose output", zflag.OptShorthand(short))

	// len(short) must be == 1
	flag := fs.ShorthandLookup(short)

	fmt.Println(flag.Name)
}
