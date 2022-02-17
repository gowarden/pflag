// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag_test

import (
	"errors"
	"fmt"
	"net"
	"os"

	"github.com/gowarden/zflag"
)

// Copyright 2020 The Go Authors. All rights reserved.
func ExampleFunc() {
	fs := zflag.NewFlagSet("ExampleFunc", zflag.ContinueOnError)
	fs.SetOutput(os.Stdout)
	var ip net.IP
	fs.Func("ip", "`IP address` to parse", func(s string) error {
		ip = net.ParseIP(s)
		if ip == nil {
			return errors.New("could not parse IP")
		}
		return nil
	})
	_ = fs.Parse([]string{"--ip", "127.0.0.1"})
	fmt.Printf("{ip: %v, loopback: %t}\n\n", ip, ip.IsLoopback())

	// 256 is not a valid IPv4 component
	_ = fs.Parse([]string{"--ip", "256.0.0.1"})
	fmt.Printf("{ip: %v, loopback: %t}\n\n", ip, ip.IsLoopback())

	// Output:
	// {ip: 127.0.0.1, loopback: true}
	//
	// invalid argument "256.0.0.1" for "--ip" flag: could not parse IP
	// Usage of ExampleFunc:
	//       --ip IP address   IP address to parse
	// {ip: <nil>, loopback: false}
}
