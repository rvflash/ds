// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package main

import (
	"flag"
	"io"
	"log"
	"os"

	"github.com/rvflash/ds/internal/mysql"
)

const (
	subCmd = iota + 1
	filePath
)

// Filled by the CI when building.
var buildVersion string

func main() {
	var (
		c1c = new(mysql.Config)
		c1f = flag.NewFlagSet(mysql.Command, flag.ExitOnError)
		w   = log.New(os.Stderr, "ds: ", 0)
		s   = "batch mode, print results using comma as the column separator, with each row on a new line"
	)
	c1f.BoolVar(&c1c.Batch, "B", false, s)
	s = "verbose mode, produce more output about what the program does"
	c1f.BoolVar(&c1c.Verbose, "v", false, s)
	s = "number of decimals to display"
	c1f.Uint64Var(&c1c.Precision, "p", mysql.DefaultPrecision, s)
	s = "number of lines to considerate by table"
	c1f.Uint64Var(&c1c.PerN, "n", mysql.DefaultPerN, s)

	var cmdName string
	if len(os.Args) > subCmd {
		cmdName = os.Args[subCmd]
	}
	switch cmdName {
	case mysql.Command:
		err := c1f.Parse(os.Args[filePath:])
		if err != nil {
			w.Fatal(err.Error())
		}
		e, err := mysql.Estimate(
			mysql.SetPrecision(c1c.Precision),
			mysql.SetPerN(c1c.PerN),
			mysql.SetBatchMode(c1c.Batch),
			mysql.SetVerbose(c1c.Verbose),
		)
		if err != nil {
			w.Fatal(err.Error())
		}
		r, err := openReader(os.Stdin, c1f.Args())
		if err != nil {
			w.Fatal(err.Error())
		}
		err = e.Run(r, os.Stdout)
		if err != nil {
			w.Fatal(err.Error())
		}
	default:
		w.Printf("version %s\n", buildVersion)
		if cmdName != "" {
			w.Fatalf("unsupported command named %q", cmdName)
		}
		w.Printf("available sub commands:\n    - %s\n", mysql.Command)
	}
}

func openReader(r io.Reader, args []string) (io.Reader, error) {
	if len(args) > 0 {
		f, err := os.Open(args[0])
		if err != nil {
			return nil, err
		}
		return f, nil
	}
	return r, nil
}
