// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag

import (
	"fmt"
	"strconv"
	"strings"
)

// -- stringToInt64 Value
type stringToInt64Value struct {
	value   *map[string]int64
	changed bool
}

func newStringToInt64Value(val map[string]int64, p *map[string]int64) *stringToInt64Value {
	ssv := new(stringToInt64Value)
	ssv.value = p
	*ssv.value = val
	return ssv
}

// Format: a=1,b=2
func (s *stringToInt64Value) Set(val string) error {
	kv := strings.SplitN(val, "=", 2)
	if len(kv) != 2 {
		return fmt.Errorf("%s must be formatted as key=value", val)
	}
	key, val := kv[0], kv[1]

	v, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return err
	}

	if !s.changed {
		*s.value = map[string]int64{}
	}

	(*s.value)[key] = v
	s.changed = true

	return nil
}

func (s *stringToInt64Value) Get() interface{} {
	return *s.value
}

func (s *stringToInt64Value) Type() string {
	return "stringToInt64"
}

func (s *stringToInt64Value) String() string {
	records := make([]string, 0, len(*s.value)>>1)
	for k, v := range *s.value {
		records = append(records, k+"="+strconv.FormatInt(v, 10))
	}

	return fmt.Sprintf("%s", records)
}

// GetStringToInt64 return the map[string]int64 value of a flag with the given name
func (f *FlagSet) GetStringToInt64(name string) (map[string]int64, error) {
	val, err := f.getFlagType(name, "stringToInt64")
	if err != nil {
		return map[string]int64{}, err
	}
	return val.(map[string]int64), nil
}

// MustGetStringToInt64 is like GetStringToInt64, but panics on error.
func (f *FlagSet) MustGetStringToInt64(name string) map[string]int64 {
	val, err := f.GetStringToInt64(name)
	if err != nil {
		panic(err)
	}
	return val
}

// StringToInt64Var defines a map[string]int64 flag with specified name, default value, and usage string.
// The argument p points to a map[string]int64 variable in which to store the values of multiple flags.
func (f *FlagSet) StringToInt64Var(p *map[string]int64, name string, value map[string]int64, usage string, opts ...Opt) {
	f.Var(newStringToInt64Value(value, p), name, usage, opts...)
}

// StringToInt64Var defines a map[string]int64 flag with specified name, default value, and usage string.
// The argument p points to a map[string]int64 variable in which to store the values of multiple flags.
func StringToInt64Var(p *map[string]int64, name string, value map[string]int64, usage string, opts ...Opt) {
	CommandLine.StringToInt64Var(p, name, value, usage, opts...)
}

// StringToInt64 defines a map[string]int64 flag with specified name, default value, and usage string.
// The return value is the address of a map[string]int64 variable that stores the values of multiple flags.
func (f *FlagSet) StringToInt64(name string, value map[string]int64, usage string, opts ...Opt) *map[string]int64 {
	var p map[string]int64
	f.StringToInt64Var(&p, name, value, usage, opts...)
	return &p
}

// StringToInt64 defines a map[string]int64 flag with specified name, default value, and usage string.
// The return value is the address of a map[string]int64 variable that stores the values of multiple flags.
func StringToInt64(name string, value map[string]int64, usage string, opts ...Opt) *map[string]int64 {
	return CommandLine.StringToInt64(name, value, usage, opts...)
}
