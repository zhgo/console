// Copyright 2014 The zhgo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package console

import (
    "strings"
)

//Request struct, containe *http.Request, it's about to more variable.
type Request struct {
    //Module name
    Module string

    //Controller name
    Controller string

    //Action name
    Action string

    //Method arguments sli.
    args []string

    //Method arguments map
    Args map[string]interface{}
}

//New request
func NewRequest(p string) *Request {
    nodes := strings.Split(p, "/")
    l := len(nodes)

    req := Request{Module: "System", Controller: "Index", Action: "Index", args: make([]string, 0), Args: make(map[string]interface{})}

    if l > 1 && nodes[1] != "" {
        req.Module = nodes[1] //strings.Title(nodes[1])
    }
    if l > 2 && nodes[2] != "" {
        req.Controller = nodes[2] //pathToMethod(nodes[2])
    }
    if l > 3 && nodes[3] != "" {
        req.Action = nodes[3] //pathToMethod(nodes[3])
    }
    if l > 4 {
        req.args = nodes[4:]
    }

    return &req
}
