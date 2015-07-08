// Copyright 2014 The zhgo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package console

import (
    "encoding/json"
    "io/ioutil"
    "strings"
    "errors"
)

type Config struct {
    // Original configuration string
    str string

    // The replacement value
    re map[string]string
}

// Load configuration file
func (c *Config) Load(p string) error {
    b, err := ioutil.ReadFile(p)
    if err != nil {
        return err
    }
    c.str = string(b)
    return nil
}

// Read from string
func (c *Config) Read(str string) error {
    if str == "" {
        return errors.New("Configuration body is empty")
    }

    c.str = str
    return nil
}

// Configuration parsing
func (c *Config) Parse(j interface{}) error {
    type variables struct {
        Variables map[string]string
    }

    r1 := variables{Variables: map[string]string{}}

    var c1 Config
    // Parsing variables node
    c1.Read(c.str)
    err := c1.parse(&r1)
    if err != nil {
        return err
    }

    // Merge
    for k, v := range r1.Variables {
        c.re[k] = v
    }

    err = c.parse(j)
    if err != nil {
        return err
    }

    return nil
}

func (c *Config) Replace(r map[string]string) *Config {
    if r == nil {
        c.re = make(map[string]string)
    } else {
        c.re = r
    }

    return c
}

// Configuration parsing
func (c *Config) parse(j interface{}) error {
    if c.str == "" {
        return errors.New("Configuration body is empty")
    }

    c.replace()

    //decode string(JSON format)
    dec := json.NewDecoder(strings.NewReader(c.str))
    err := dec.Decode(j)
    if err != nil {
        return err
    }

    return nil
}

//Replace placeholer({})
func (c *Config) replace() {
    if len(c.re) == 0 {
        return
    }

    for k, v := range c.re {
        //c.str = regexp.MustCompile(k).ReplaceAllLiteralString(c.str, v)
        c.str = strings.Replace(c.str, k, v, -1)
    }
}

// New Config
func NewConfig(p string) *Config {
    c := Config{re: make(map[string]string)}
    c.Load(p)
    return &c
}