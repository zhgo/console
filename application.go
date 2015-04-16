// Copyright 2014 The zhgo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package console

import (
    "github.com/zhgo/db"
    "github.com/zhgo/kernel"
    "github.com/zhgo/config"
    "os"
    "fmt"
)

// Application struct
type Application struct {
    // Environment 0:development 1:testing 2:staging 3:production
    Environment int8

    // module list
    Modules map[string]Module

    // Environments of operation system.
    Env map[string]string

    // Path of operation system.
    Path []string

    //
    Programs map[string]string

    //
    Services map[string][]string

    //
    AutoRuns []string
}

// Module struct
type Module struct {
    // module name
    Name string

    // key of DSN
    DB db.Server
}

// App
var App Application

// Init
func (app *Application) Init(path string) {
    // Load config file
    r := map[string]string{"{WorkingDir}": kernel.WorkingDir}
    config.LoadJSONFileAdv(app, path, r)

    // Default module
    if app.Modules == nil {
        app.Modules = make(map[string]Module)
    }

    //Env
    for k, v := range app.Env {
        Setenv(k, v)
    }

    //Path
    p := os.Getenv("Path")
    for _, v := range app.Path {
        p = fmt.Sprintf("%s%c%s", p, os.PathListSeparator, v)
    }
    Setenv("PATH", p)
}

// Load
func (app *Application) Load(p, args string) {

}

// Start
func (app *Application) Start() {

}
