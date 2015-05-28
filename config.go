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

// Configuration content type
const (
    ConfigFile = iota
    ConfigString
)

type Config struct {
    // Original configuration string
    str string

    // The replacement value
    rep map[string]string
}

// Read from string
func (c *Config) Read(str string) error {
    if str == "" {
        return errors.New("Configuration body is empty")
    }

    c.str = str
    return nil
}

// Load configuration file
func (c *Config) Load(path string) error {
    b, err := ioutil.ReadFile(path)
    if err != nil {
        return err
    }
    c.str = string(b)
    return nil
}

// Configuration parsing
func (c *Config) parse(j interface{}) error {
    if c.str == "" {
        return errors.New("Configuration body is empty")
    }

    c.str = c.replace(c.str, c.rep)

    //decode string(JSON format)
    dec := json.NewDecoder(strings.NewReader(c.str))
    err := dec.Decode(j)
    if err != nil {
        return err
    }

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
        c.rep[k] = v
    }

    err = c.parse(j)
    if err != nil {
        return err
    }

    return nil
}

//Replace by map
func (c *Config) replace(t string, replaces map[string]string) string {
    if len(replaces) == 0 {
        return t
    }

    for k, v := range replaces {
        //t = regexp.MustCompile(k).ReplaceAllLiteralString(t, v)
        t = strings.Replace(t, k, v, -1)
    }

    return t
}

// New Config
func NewConfig(typ uint8, str string, rep map[string]string) *Config {
    var c Config
    if typ == ConfigFile {
        c.Load(str)
    } else {
        c.Read(str)
    }
    if rep == nil {
        c.rep = make(map[string]string)
    } else {
        c.rep = rep
    }
    return &c
}