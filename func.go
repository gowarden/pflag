// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag

// -- int Value
type funcValue func(string) error

func newFuncValue(fn func(string) error) *funcValue {
	funcVal := funcValue(fn)
	return &funcVal
}

func (i *funcValue) Set(s string) error {
	return (*i)(s)
}

func (i *funcValue) Type() string {
	return "string"
}

func (i *funcValue) String() string { return "" }

// Func defines a flag with specified name, and usage string.
// Each time the flag is seen, fn is called with the value of the flag.
// If fn returns a non-nil error, it will be treated as a flag value parsing error.
func (f *FlagSet) Func(name string, usage string, fn func(string) error, opts ...Opt) {
	f.Var(newFuncValue(fn), name, usage, opts...)
}

// Func defines a flag with specified name, and usage string.
// Each time the flag is seen, fn is called with the value of the flag.
// If fn returns a non-nil error, it will be treated as a flag value parsing error.
func Func(name string, usage string, fn func(string) error, opts ...Opt) {
	CommandLine.Func(name, usage, fn, opts...)
}
