// Copyright 2014 The zhgo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package console

import (
	"bufio"
	"fmt"
	"github.com/zhgo/config"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// typ: 0:service, 1:command
const (
	ProgramSRV = iota
	ProgramCMD
)

// Console struct
type Console struct {
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
	cmdText string

	cmds []*exec.Cmd
}

func (c *Console) Init(path string) {
	type variables struct {
		Variables map[string]string
	}
	var r variables
	err := config.LoadJSONFile(&r, path, nil)
	if err != nil {
		log.Printf("%s\n", err)
	}

	//r.Variables["{basePath}"] = WorkingDir
	//r.Variables["{SystemRoot}"] = console.Getenv("SystemRoot")
	//r.Variables["{UserProfile}"] = console.Getenv("USERPROFILE")

	err = config.LoadJSONFile(c, path, r.Variables)
	if err != nil {
		log.Printf("%s\n", err)
	}

	//Env
	for k, v := range c.Env {
		Setenv(k, v)
	}

	//Path
	p := os.Getenv("Path")
	for _, v := range c.Path {
		p = fmt.Sprintf("%s%c%s", p, os.PathListSeparator, v)
	}
	Setenv("PATH", p)

	c.cmds = make([]*exec.Cmd, 0)
}

func (c *Console) AutoRun() {
	for _, v := range c.AutoRuns {
		c.ServiceItem("start", v)
	}
}

func (c *Console) Start(perfix string) {
consoleloop:
	for {
		path, err := os.Getwd()
		if err != nil {
			fmt.Printf("%s\n", err)
			break consoleloop
		}

		fmt.Printf("[%s@%s] ", perfix, filepath.Base(path))

		reader := bufio.NewReader(os.Stdin)
		strBytes, _, err := reader.ReadLine()
		if err != nil {
			fmt.Printf("%s\n", err)
		}

		c.cmdText = string(strBytes)
		args := ParseText(c.cmdText)
		c.RunProgram(args)
	}
}

func (c *Console) Exit() {
	//for _, cmd := range c.cmds {
	//	cmd.Process.Kill()
	//}

	for _, v := range c.AutoRuns {
		c.ServiceItem("stop", v)
	}

	os.Exit(0)
}

func (c *Console) RunProgram(args []string) {
	if len(args) == 0 {
		return
	}

	//run Programs
	cmd, s := c.Programs[args[0]]
	if s {
		sli := ParseText(cmd)
		args = append(sli, args[1:]...)
		c.RunProgram(args)
	} else {
		switch args[0] {
		case "new":
			if runtime.GOOS == "windows" {
				c.RunProgram(append([]string{"cmd.exe", "/C", "start"}, args[1:]...))
			} else {
				//FIXME: support darwin, freebsd, linux
			}
		case "async":
			go c.RunProgram(args[1:])
		case "cd":
			Chdir(args[1])
		case "cls", "del", "deltree", "dir", "path", "set": //inner command
			if runtime.GOOS == "windows" {
				c.RunProgram(append([]string{"cmd.exe", "/C"}, args...))
			} else {
				//FIXME: support darwin, freebsd, linux
			}
		case "exit", "q", "quit":
			c.Exit()
		case "service":
			if len(args) > 2 {
				c.Service(args[1], args[2])
			}
		default:
			//debug.Dump(args)
			c.ExecuteCMD(args[0], args[1:]...)
		}
	}
}

func (c *Console) ExecuteCMD(path string, args ...string) *exec.Cmd {
	cmd := exec.Command(path, args...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		//fmt.Printf("%s\n", err)
	}

	return cmd
}

func (c *Console) ExecuteSRV(path string, args ...string) *exec.Cmd {
	cmd := exec.Command(path, args...)

	fmt.Printf("[%s] %v\n", path, args)

	err := cmd.Start()
	if err != nil {
		fmt.Printf("%s\n", err)
	} else {
		fmt.Println(cmd.Process.Pid)
	}

	return cmd
}

func (c *Console) Service(typ string, name string) {
	if name == "all" {
		for k, _ := range c.Services {
			c.ServiceItem(typ, k)
		}
	} else {
		c.ServiceItem(typ, name)
	}
}

func (c *Console) ServiceItem(typ string, name string) {
	services, s := c.Services[name]
	if s == false {
		log.Printf("Service not existing: %s\n", name)
		return
	}

	if typ == "start" {
		for _, p := range services {
			args := ParseText(p)
			if args[0] == "start" {
				cmd := c.ExecuteSRV(args[1], args[2:]...)
				c.cmds = append(c.cmds, cmd)
			}
			if args[0] == "run" {
				c.RunProgram(args[1:])
			}
		}
	} else if typ == "stop" {
		for _, p := range services {
			args := ParseText(p)
			if args[0] == "start" {
				if runtime.GOOS == "windows" {
					c.ExecuteCMD("taskkill", "/f", "/im", args[1])
				} else {
					//FIXME: support darwin, freebsd, linux

				}
			}
		}
	}
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
