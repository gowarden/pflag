// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag

import (
	"fmt"
	"strconv"
	"strings"
)

// -- stringToString Value
type stringToStringValue struct {
	value   *map[string]string
	changed bool
}

func newStringToStringValue(val map[string]string, p *map[string]string) *stringToStringValue {
	ssv := new(stringToStringValue)
	ssv.value = p
	*ssv.value = val
	return ssv
}

func (s *stringToStringValue) Set(val string) error {
	kv := strings.SplitN(val, "=", 2)
	if len(kv) != 2 {
		return fmt.Errorf("%s must be formatted as key=value", val)
	}
	key, val := kv[0], kv[1]

	if !s.changed {
		*s.value = map[string]string{}
	}

	(*s.value)[key] = val
	s.changed = true

	return nil
}

func (s *stringToStringValue) Get() interface{} {
	return *s.value
}

func (s *stringToStringValue) Type() string {
	return "stringToString"
}

func (s *stringToStringValue) String() string {
	records := make([]string, 0, len(*s.value)>>1)
	for k, v := range *s.value {
		records = append(records, k+"="+strconv.Quote(v))
	}

	return fmt.Sprintf("%s", records)
}

// GetStringToString return the map[string]string value of a flag with the given name
func (f *FlagSet) GetStringToString(name string) (map[string]string, error) {
	val, err := f.getFlagType(name, "stringToString")
	if err != nil {
		return map[string]string{}, err
	}
	return val.(map[string]string), nil
}

// MustGetStringToString is like GetStringToString, but panics on error.
func (f *FlagSet) MustGetStringToString(name string) map[string]string {
	val, err := f.GetStringToString(name)
	if err != nil {
		panic(err)
	}
	return val
}

// StringToStringVar defines a map[string]string flag with specified name, default value, and usage string.
// The argument p points to a map[string]string variable in which to store the values of multiple flags.
func (f *FlagSet) StringToStringVar(p *map[string]string, name string, value map[string]string, usage string, opts ...Opt) {
	f.Var(newStringToStringValue(value, p), name, usage, opts...)
}

// StringToStringVar defines a map[string]string flag with specified name, default value, and usage string.
// The argument p points to a map[string]string variable in which to store the values of multiple flags.
func StringToStringVar(p *map[string]string, name string, value map[string]string, usage string, opts ...Opt) {
	CommandLine.StringToStringVar(p, name, value, usage, opts...)
}

// StringToString defines a map[string]string flag with specified name, default value, and usage string.
// The return value is the address of a map[string]string variable that stores the values of multiple flags.
func (f *FlagSet) StringToString(name string, value map[string]string, usage string, opts ...Opt) *map[string]string {
	var p map[string]string
	f.StringToStringVar(&p, name, value, usage, opts...)
	return &p
}

// StringToString defines a map[string]string flag with specified name, default value, and usage string.
// The return value is the address of a map[string]string variable that stores the values of multiple flags.
func StringToString(name string, value map[string]string, usage string, opts ...Opt) *map[string]string {
	return CommandLine.StringToString(name, value, usage, opts...)
}
