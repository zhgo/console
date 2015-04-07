// Copyright 2014 The zhgo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package console

import (
    "reflect"
)

//Registered controllers
var controllers map[string]reflect.Value = make(map[string]reflect.Value)

//controller struct
type Controller struct {
    //Request
    Request *Request
}

//new router register
func NewController(module string, c interface{}) {
    value := reflect.ValueOf(c)
    controllers[value.Elem().Type().Name()] = value
}
