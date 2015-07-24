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
