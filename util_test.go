// Copyright 2014 The zhgo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package console

import (
	"testing"
)

func TestUnderscoreToCamelcase(t *testing.T) {
	path := "browse_by_set"
	method := UnderscoreToCamelcase(path)
	if method != "BrowseBySet" {
		t.Error("pathToMethod failure")
	}
}

func TestCamelcaseToUnderscore(t *testing.T) {
	method := "BrowseBySet"
	path := CamelcaseToUnderscore(method)
	if path != "browse_by_set" {
		t.Errorf("methodToPath failure: %#v", path)
	}
}

func TestParseText(t *testing.T) {
	cmd := "ls.exe -l"
	sli := ParseText(cmd)
	if len(sli) != 2 || sli[0] != "ls.exe" || sli[1] != "-l" {
		t.Errorf("ParseText failure: %#v", cmd)
	}

	cmd = "cd {homePath}/gocode/src"
	sli = ParseText(cmd)
	if len(sli) != 2 || sli[0] != "cd" || sli[1] != "{homePath}/gocode/src" {
		t.Errorf("ParseText failure: %#v", cmd)
	}

	cmd = "run cd {basePath}/usr/local/apache24"
	sli = ParseText(cmd)
	if len(sli) != 3 || sli[0] != "run" || sli[1] != "cd" || sli[2] != "{basePath}/usr/local/apache24" {
		t.Errorf("ParseText failure: %#v", cmd)
	}

	cmd = "srv httpd.exe"
	sli = ParseText(cmd)
	if len(sli) != 2 || sli[0] != "srv" || sli[1] != "httpd.exe" {
		t.Errorf("ParseText failure: %#v", cmd)
	}

	cmd = "async 'D:/Program Files/Sublime Text 3/sublime_text.exe'"
	sli = ParseText(cmd)
	if len(sli) != 2 {
		t.Errorf("ParseText failure: %#v", cmd)
	}
}
