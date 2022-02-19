// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag

import (
	"fmt"
	"strconv"
	"strings"
)

// -- int32Slice Value
type int32SliceValue struct {
	value   *[]int32
	changed bool
}

func newInt32SliceValue(val []int32, p *[]int32) *int32SliceValue {
	isv := new(int32SliceValue)
	isv.value = p
	*isv.value = val
	return isv
}

func (s *int32SliceValue) Set(val string) error {
	val = strings.TrimSpace(val)
	temp64, err := strconv.ParseInt(val, 0, 32)
	if err != nil {
		return err
	}

	if !s.changed {
		*s.value = []int32{}
	}
	*s.value = append(*s.value, int32(temp64))
	s.changed = true

	return nil
}

func (s *int32SliceValue) Get() interface{} {
	return *s.value
}

func (s *int32SliceValue) Type() string {
	return "int32Slice"
}

func (s *int32SliceValue) String() string {
	if s.value == nil {
		return "[]"
	}

	return fmt.Sprintf("%d", *s.value)
}

func (s *int32SliceValue) fromString(val string) (int32, error) {
	t64, err := strconv.ParseInt(val, 0, 32)
	if err != nil {
		return 0, err
	}
	return int32(t64), nil
}

func (s *int32SliceValue) toString(val int32) string {
	return fmt.Sprintf("%d", val)
}

func (s *int32SliceValue) Append(val string) error {
	i, err := s.fromString(val)
	if err != nil {
		return err
	}
	*s.value = append(*s.value, i)
	return nil
}

func (s *int32SliceValue) Replace(val []string) error {
	out := make([]int32, len(val))
	for i, d := range val {
		var err error
		out[i], err = s.fromString(d)
		if err != nil {
			return err
		}
	}
	*s.value = out
	return nil
}

func (s *int32SliceValue) GetSlice() []string {
	out := make([]string, len(*s.value))
	for i, d := range *s.value {
		out[i] = s.toString(d)
	}
	return out
}

// GetInt32Slice return the []int32 value of a flag with the given name
func (f *FlagSet) GetInt32Slice(name string) ([]int32, error) {
	val, err := f.getFlagType(name, "int32Slice")
	if err != nil {
		return []int32{}, err
	}
	return val.([]int32), nil
}

// MustGetInt32Slice is like GetInt32Slice, but panics on error.
func (f *FlagSet) MustGetInt32Slice(name string) []int32 {
	val, err := f.GetInt32Slice(name)
	if err != nil {
		panic(err)
	}
	return val
}

// Int32SliceVar defines a []int32 flag with specified name, default value, and usage string.
// The argument p points to a []int32 variable in which to store the value of the flag.
func (f *FlagSet) Int32SliceVar(p *[]int32, name string, value []int32, usage string, opts ...Opt) {
	f.Var(newInt32SliceValue(value, p), name, usage, opts...)
}

// Int32SliceVar defines a []int32 flag with specified name, default value, and usage string.
// The argument p points to a []int32 variable in which to store the value of the flag.
func Int32SliceVar(p *[]int32, name string, value []int32, usage string, opts ...Opt) {
	CommandLine.Int32SliceVar(p, name, value, usage, opts...)
}

// Int32Slice defines a []int32 flag with specified name, default value, and usage string.
// The return value is the address of a []int32 variable that stores the value of the flag.
func (f *FlagSet) Int32Slice(name string, value []int32, usage string, opts ...Opt) *[]int32 {
	var p []int32
	f.Int32SliceVar(&p, name, value, usage, opts...)
	return &p
}

// Int32Slice defines a []int32 flag with specified name, default value, and usage string.
// The return value is the address of a []int32 variable that stores the value of the flag.
func Int32Slice(name string, value []int32, usage string, opts ...Opt) *[]int32 {
	return CommandLine.Int32Slice(name, value, usage, opts...)
}
