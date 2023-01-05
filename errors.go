// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag

import (
	"fmt"
	"strings"
)

func getFlagWithDashes(name string) string {
	dash := "--"
	if len(name) == 1 {
		dash = "-"
	}

	return dash + name
}

type UnknownFlagError struct {
	name string
}

var _ error = (*UnknownFlagError)(nil)

func NewUnknownFlagError(name string) error {
	return UnknownFlagError{name: name}
}

func (e UnknownFlagError) Error() string {
	return fmt.Sprintf("unknown flag: %s", getFlagWithDashes(e.name))
}

type MissingFlagsError []string

var _ error = (*MissingFlagsError)(nil)

func (e *MissingFlagsError) AddMissingFlag(f *Flag) {
	*e = append(*e, getFlagWithDashes(f.Name))
}

func (e MissingFlagsError) Error() string {
	flagNames := make([]string, 0, len(e))
	for _, s := range e {
		flagNames = append(flagNames, fmt.Sprintf("%q", s))
	}

	return fmt.Sprintf(`required flag(s) %s not set`, strings.Join(flagNames, `, `))
}

type InvalidArgumentError struct {
	flagName string
	value    interface{}
	err      error
}

var _ error = (*InvalidArgumentError)(nil)

func NewInvalidArgumentError(err error, f *Flag, value interface{}) error {
	var flagName string
	if f.Shorthand != 0 && f.ShorthandDeprecated == "" {
		flagName = fmt.Sprintf("-%c", f.Shorthand)
		if !f.ShorthandOnly {
			flagName = fmt.Sprintf("%s, --%s", flagName, f.Name)
		}
	} else {
		flagName = getFlagWithDashes(f.Name)
	}

	return InvalidArgumentError{
		flagName: flagName,
		value:    value,
		err:      err,
	}
}

func (e InvalidArgumentError) Error() string {
	return fmt.Sprintf("invalid argument %q for %q flag: %s", e.value, e.flagName, e.err)
}

func (e InvalidArgumentError) Unwrap() error {
	return e.err
}
