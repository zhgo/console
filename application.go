// Copyright 2014 The zhgo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package console

import (
    "os"
    "fmt"
    "path/filepath"
    "bufio"
    "runtime"
    "os/exec"
    "log"
)

// Application struct
type Application struct {
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

    //
    promptPerfix string

    //
    CmdList []*exec.Cmd
}

// Init
func (app *Application) Init(path string) {
    // Load config file
    r := map[string]string{"{WorkingDir}": WorkingDir}
    NewConfig(path).Replace(r).Parse(app)

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

    // CmdList
    app.CmdList = make([]*exec.Cmd,0)
}

// Load
func (app *Application) Load(promptPerfix string) {
    app.promptPerfix = promptPerfix

    // Auto run
    for _, v := range app.AutoRuns {
        app.RunSRV(v)
    }
}

// Start
func (app *Application) Start() {
    consoleloop:
    for {
        path, err := os.Getwd()
        if err != nil {
            fmt.Printf("%s\n", err)
            break consoleloop
        }

        fmt.Printf("[%s@%s] ", app.promptPerfix, filepath.Base(path))

        reader := bufio.NewReader(os.Stdin)
        strBytes, _, err := reader.ReadLine()
        if err != nil {
            fmt.Printf("%s\n", err)
        }

        cmdText := string(strBytes)
        args := ParseText(cmdText)
        app.Run(args)
    }
}

func (app *Application) Run(args []string) {
    if len(args) == 0 {
        return
    }

    //run Programs
    cmd, s := app.Programs[args[0]]
    if s {
        sli := ParseText(cmd)
        args = append(sli, args[1:]...)
        app.Run(args)
        return
    }

    switch args[0] {
        // Change directory
        case "cd":
        Chdir(args[1])

        // Inner commands
        case "cls", "del", "deltree", "dir", "path", "set":
        if runtime.GOOS == "windows" {
            app.Run(append([]string{"cmd.exe", "/C"}, args...))
        } else {
            //FIXME: support darwin, freebsd, linux
        }

        // open in new window
        case "new":
        if runtime.GOOS == "windows" {
            app.Run(append([]string{"cmd.exe", "/C", "start"}, args[1:]...))
        } else {
            //FIXME: support darwin, freebsd, linux
        }

        // Asynchronous run program
        case "async":
        go app.Run(args[1:])

        // Run as service
        case "srv":
        app.ExecuteSRV(args[1], args[2:]...)

        // Run as Program
        case "run":
        app.Run(args[1:])

        // Exit
        case "exit", "q", "quit":
        app.Exit()

        // Default run as application
        default:
        app.ExecuteCMD(args[0], args[1:]...)
    }
}

func (app *Application) RunSRV(name string) {
    services, s := app.Services[name]
    if s == false {
        log.Printf("Service not existing: %s\n", name)
        return
    }

    for _, p := range services {
        args := ParseText(p)
        app.Run(args)
    }
}

func (app *Application) ExecuteCMD(path string, args ...string) {
    cmd := exec.Command(path, args...)
    // app.CmdList = append(app.CmdList, cmd)

    fmt.Printf("[%s] %v\n", path, args)

    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

    err := cmd.Run()
    if err != nil {
        fmt.Printf("%s\n", err)
    }

    fmt.Printf("\n")
}

func (app *Application) ExecuteSRV(path string, args ...string) {
    cmd := exec.Command(path, args...)
    // app.CmdList = append(app.CmdList, cmd)

    fmt.Printf("[%s] %v\n", path, args)

    err := cmd.Start()
    if err != nil {
        fmt.Printf("%s\n", err)
    } else {
        fmt.Printf("%d\n", cmd.Process.Pid)
    }

    fmt.Printf("\n")
}

func (app *Application) Exit() {
    app.RunSRV("stopall")

    os.Exit(0)
}
