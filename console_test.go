// Copyright 2014 The zhgo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package console

import (
	"testing"
)

var cmd string
var args []string

func TestParseText(t *testing.T) {
	cmd = "go build github.com/liudng/dogo"
	args = ParseText(cmd)
	if len(args) != 3 {
		t.Fatalf("ParseText failure: \n%#v\n%#v\n", cmd, args)
	}
	if args[0] != "go" || args[1] != "build" || args[2] != "github.com/liudng/dogo" {
		t.Fatalf("ParseText failure: \n%#v\n%#v\n", cmd, args)
	}

	cmd = "\"D:/Program Files (x86)/liteide/bin/liteide.exe\" -cli"
	args = ParseText(cmd)
	if len(args) != 2 {
		t.Fatalf("ParseText failure: \n%#v\n%#v\n", cmd, args)
	}
	if args[0] != "D:/Program Files (x86)/liteide/bin/liteide.exe" || args[1] != "-cli" {
		t.Fatalf("ParseText failure: \n%#v\n%#v\n", cmd, args)
	}
}
