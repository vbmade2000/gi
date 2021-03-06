package main

import (
	"flag"
	"fmt"
	"github.com/go-interpreter/gi/pkg/compiler"
	"os"
	"path"
	"path/filepath"
)

var ProgramName string = path.Base(os.Args[0])

type GIConfig struct {
	Quiet          bool
	Verbose        bool
	VerboseVerbose bool
	RawLua         bool
	PreludePath    string

	preludeFiles []string
}

// call DefineFlags before myflags.Parse()
func (c *GIConfig) DefineFlags(fs *flag.FlagSet) {
	fs.BoolVar(&c.Quiet, "q", false, "don't show banner on startup")
	fs.BoolVar(&c.Verbose, "v", false, "show debug prints")
	fs.BoolVar(&c.VerboseVerbose, "vv", false, "show even more verbose debug prints")
	fs.BoolVar(&c.RawLua, "raw", false, "skip all translation, type raw Lua to LuaJIT with our prelude installed")
	fs.StringVar(&c.PreludePath, "prelude", "", "path to the prelude directory. All .lua files are sourced before startup from this directory. Default is to to read from 'GOINTERP_PRELUD_DIR' env var. -prelude overrides this.")
}

var defaultPreludePathParts = []string{
	"src",
	"github.com",
	"go-interpreter",
	"gi",
	"pkg",
	"compiler"}

// call c.ValidateConfig() after myflags.Parse()
func (c *GIConfig) ValidateConfig() error {

	if c.PreludePath == "" {
		dir := os.Getenv("GOINTERP_PRELUDE_DIR")
		if dir != "" {
			c.PreludePath = dir
		} else {
			// try hard... try $GOPATH/src/github.com/go-interpreter/gi/pkg/compiler
			// by default.
			gopath := os.Getenv("GOPATH")
			if gopath == "" {
				// try $HOME/go
				home := os.Getenv("HOME")
				proposed := filepath.Join(home, "go")
				if !DirExists(home) || !DirExists(proposed) {
					return preludeError()
				}
				gopath = proposed
			}

			c.PreludePath = filepath.Join(append([]string{gopath}, defaultPreludePathParts...)...)
		}
	}
	files, err := compiler.FetchPrelude(c.PreludePath)
	if err != nil {
		return err
	}
	c.preludeFiles = files
	return nil
}

func preludeError() error {
	return fmt.Errorf("setenv GOINTERP_PRELUDE_DIR to point to your prelude dir. This is typically $GOPATH/src/github.com/go-interpreter/gi/pkg/compiler but GOINTERP_PRELUDE_DIR was not set and -prelude was not specified.")
}
