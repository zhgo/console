// Copyright 2014 The zhgo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package console

import (
    "log"
    "os"
    "runtime"
    "strings"
)

// Working Directory
var WorkingDir string = workingDir()

//Get working directory.
func workingDir() string {
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

func Setenv(key, value string) {
    err := os.Setenv(key, value)
    if err != nil {
        log.Printf("%s\n", err)
    }
}

func Chdir(path string) {
    err := os.Chdir(path)
    if err != nil {
        log.Printf("%s\n", err)
    }
}

func Getenv(key string) string {
    value := os.Getenv(key)
    if runtime.GOOS == "windows" {
        value = strings.Replace(value, "\\", "/", -1)
    }

    return value
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

