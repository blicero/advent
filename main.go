// /home/krylon/go/src/github.com/blicero/advent/main.go
// -*- mode: go; coding: utf-8; -*-
// Created on 19. 10. 2026 by Benjamin Walkenhorst
// (c) 2026 Benjamin Walkenhorst
// Time-stamp: <2026-10-19 08:40:24 krylon>

package main

import (
	"fmt"
	"os"

	"github.com/blicero/advent/common"
	"github.com/blicero/advent/prompt"
)

func main() {
	var (
		err error
		sh  *prompt.Shell
	)

	defer fmt.Println("Toodles!")

	if err = common.InitApp(); err != nil {
		fmt.Fprintf(
			os.Stderr,
			"Failed to initialize AppDir %s: %s\n",
			common.BaseDir,
			err.Error())
		os.Exit(1)
	} else if sh, err = prompt.Create(nil); err != nil {
		fmt.Fprintf(
			os.Stderr,
			"Cannot create Shell: %s\n",
			err.Error())
		os.Exit(1)
	}

	sh.Run()
} // func main()
