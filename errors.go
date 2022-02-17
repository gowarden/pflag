// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag

import "fmt"

type errUnknownFlag struct {
	name string
}

func NewUnknownFlagError(name string) error {
	return errUnknownFlag{name: name}
}

func (e errUnknownFlag) Error() string {
	return fmt.Sprintf("unknown flag: %s", getFlagWithDashes(e.name))
}
