// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag

import (
	"fmt"
)

type FlagUsageFormatter func(*Flag) (string, string)

func defaultUsageFormatter(flag *Flag) (string, string) {
	left := "  "
	if flag.Shorthand != 0 && flag.ShorthandDeprecated == "" {
		left += fmt.Sprintf("-%c", flag.Shorthand)
		if !flag.ShorthandOnly {
			left += ", "
		}
	} else {
		left += "    "
	}
	left += "--"
	if _, ok := flag.Value.(boolFlag); ok {
		left += "[no-]"
	}
	left += flag.Name

	varname, usage := UnquoteUsage(flag)
	if varname != "" {
		left += " " + varname
	}

	right := usage
	if !flag.DisablePrintDefault && !flag.defaultIsZeroValue() {
		if v, ok := flag.Value.(Typed); ok && v.Type() == "string" {
			right += fmt.Sprintf(" (default %q)", flag.DefValue)
		} else {
			right += fmt.Sprintf(" (default %s)", flag.DefValue)
		}
	}
	if len(flag.Deprecated) != 0 {
		right += fmt.Sprintf(" (DEPRECATED: %s)", flag.Deprecated)
	}

	return left, right
}
