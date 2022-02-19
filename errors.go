// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag

import "fmt"

type UnknownFlagError struct {
	name string
}

func NewUnknownFlagError(name string) error {
	return UnknownFlagError{name: name}
}

func (e UnknownFlagError) Error() string {
	return fmt.Sprintf("unknown flag: %s", getFlagWithDashes(e.name))
}
