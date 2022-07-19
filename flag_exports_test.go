// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag

func SetExitFunc(fn func(code int)) {
	exitFn = fn
}
