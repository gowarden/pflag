// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag_test

import (
	"fmt"

	"github.com/zulucmd/zflag/v2"
)

func ExampleShorthandLookup() {
	name := "verbose"
	short := 'v'

	zflag.Bool(name, false, "verbose output", zflag.OptShorthand(short))

	// len(short) must be == 1
	flag := zflag.ShorthandLookup(short)

	fmt.Println(flag.Name)

	// Output:
	// verbose
}

func ExampleFlagSet_ShorthandLookup() {
	name := "verbose"
	short := 'v'

	fs := zflag.NewFlagSet("Example", zflag.ContinueOnError)
	fs.Bool(name, false, "verbose output", zflag.OptShorthand(short))

	// len(short) must be == 1
	flag := fs.ShorthandLookup(short)

	fmt.Println(flag.Name)

	// Output:
	// verbose
}

func ExampleFlag_Required() {
	fs := zflag.NewFlagSet("Example", zflag.ContinueOnError)
	fs.Bool("required", false, "flag must be set", zflag.OptRequired())
	err := fs.Parse([]string{})
	fmt.Println(err)

	// Output:
	// required flag(s) "--required" not set
}
