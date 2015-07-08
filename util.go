// Copyright 2014 The zhgo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package console

import (
    "log"
    "os"
    "runtime"
    "strings"
    "flag"
    "fmt"
)

// Working Directory
var WorkingDir string = Getwd()

// Console Application
var App Application

//Wrap os.Getwd()
func Getwd() string {
    w, err := os.Getwd()
    if err != nil {
        log.Fatal(err)
    }

    //replace backslash(\) to slash(/) on windows platform
    if runtime.GOOS == "windows" {
        w = strings.Replace(w, "\\", "/", -1)
    }

    return w
}

func Getenv(key string) string {
    value := os.Getenv(key)
    if runtime.GOOS == "windows" {
        value = strings.Replace(value, "\\", "/", -1)
    }

    return value
}

func Setenv(key, value string) error {
    err := os.Setenv(key, value)
    if err != nil {
        log.Printf("%s\n", err)
        return err
    }

    return nil
}

func Chdir(path string) error {
    err := os.Chdir(path)
    if err != nil {
        log.Printf("%s\n", err)
        return err
    }

    return nil
}

func ParseText(txt string) []string {
    //support space in path
    args := make([]string, 0)
    node := ""
    colon := false
    for _, c := range txt {
        t := string(c)
        //log.Printf("%#v\n", t)
        if t == "'" {
            if colon == true {
                colon = false
            } else {
                colon = true
            }

            continue
        }

        if (t == " " && colon == true) || t != " " {
            node += t
        } else {
            args = append(args, node)
            node = ""
        }
    }

    if node != "" {
        args = append(args, node)
    }

    //args := strings.Split(txt, " ")

    return args
}

// Console parameters
func Arguments(app string) (string, string) {
    var c, h, p string
    flag.StringVar(&c, "c", WorkingDir+"/example.json", "Usage: mplus -c=/path/to/example.json")
    flag.StringVar(&h, "h", "nil", "Usage: example -h")
    flag.StringVar(&p, "p", "", "Usage: example -p=Passport/User/Login&id=1")
    flag.Parse()

    if h != "nil" {
        fmt.Println(
        fmt.Sprintf("Usage: %s [OPTION]...", app),
        "{example} is the name of the application, you can change in a real environment.",
        "",
        "  -c  The path of the configuration file.",
        "  -h  Display this help and exit.",
        "  -p  Console application action path. Separated by a slash.")
        os.Exit(0)
    }

    return c, p
}