// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag

import (
	"fmt"
)

type Opt func(f *Flag) error

func applyFlagOptions(f *Flag, options ...Opt) error {
	for _, option := range options {
		if err := option(f); err != nil {
			return err
		}
	}
	return nil
}

// OptShorthand one-letter abbreviated flag
func OptShorthand(shorthand rune) Opt {
	return func(f *Flag) error {
		f.Shorthand = shorthand
		return nil
	}
}

// OptShorthandStr one-letter abbreviated flag
func OptShorthandStr(shorthand string) Opt {
	r, err := shorthandStrToRune(shorthand)
	if err != nil {
		panic(err)
	}

	return OptShorthand(r)
}

// OptShorthandOnly If the user set only the shorthand
func OptShorthandOnly() Opt {
	return func(f *Flag) error {
		f.ShorthandOnly = true
		return nil
	}
}

// OptUsageType flag type displayed in the help message
func OptUsageType(usageType string) Opt {
	return func(f *Flag) error {
		f.UsageType = usageType
		return nil
	}
}

// OptDisableUnquoteUsage disable unquoting and extraction of type from usage
func OptDisableUnquoteUsage() Opt {
	return func(f *Flag) error {
		f.DisableUnquoteUsage = true
		return nil
	}
}

// OptDisablePrintDefault toggle printing of the default value in usage message
func OptDisablePrintDefault() Opt {
	return func(f *Flag) error {
		f.DisablePrintDefault = true
		return nil
	}
}

// OptDefValue default value (as text); for usage message
func OptDefValue(defValue string) Opt {
	return func(f *Flag) error {
		f.DefValue = defValue
		return nil
	}
}

// OptNoOptDefVal default value (as text); if the flag is on the command line without any options
func OptNoOptDefVal(noOptDefVal string) Opt {
	return func(f *Flag) error {
		f.NoOptDefVal = noOptDefVal
		return nil
	}
}

// OptDeprecated indicated that a flag is deprecated in your program. It will
// continue to function but will not show up in help or usage messages. Using
// this flag will also print the given usageMessage.
func OptDeprecated(msg string) Opt {
	return func(f *Flag) error {
		if msg == "" {
			return fmt.Errorf("deprecated message for flag %q must be set", f.Name)
		}

		f.Deprecated = msg
		return OptHidden()(f)
	}
}

// OptHidden used by zulu.Command to allow flags to be hidden from help/usage text
func OptHidden() Opt {
	return func(f *Flag) error {
		f.Hidden = true
		return nil
	}
}

// OptShorthandDeprecated If the shorthand of this flag is deprecated, this string is the new or now thing to use
func OptShorthandDeprecated(msg string) Opt {
	return func(f *Flag) error {
		if msg == "" {
			return fmt.Errorf("shorthand deprecated message for flag %q must be set", f.Name)
		}

		f.ShorthandDeprecated = msg
		return nil
	}

}

// OptGroup flag group
func OptGroup(group string) Opt {
	return func(f *Flag) error {
		f.Group = group
		return nil
	}
}

// OptAnnotation Use it to annotate this specific flag for your application
func OptAnnotation(key string, value []string) Opt {
	return func(f *Flag) error {
		return f.SetAnnotation(key, value)
	}
}
