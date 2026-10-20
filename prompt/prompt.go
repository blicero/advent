// /home/krylon/go/src/github.com/blicero/advent/prompt/prompt.go
// -*- mode: go; coding: utf-8; -*-
// Created on 18. 10. 2026 by Benjamin Walkenhorst
// (c) 2026 Benjamin Walkenhorst
// Time-stamp: <2026-10-19 10:11:24 krylon>

// Package prompt handles the interaction with the user.
// Like, ewww, people.
package prompt

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/blicero/advent/common"
	"github.com/blicero/advent/logdomain"
	"github.com/blicero/advent/model"
	gprompt "github.com/c-bata/go-prompt"
)

const prefix = "~~> "

// Shell wraps the interaction with the terminal.
type Shell struct {
	log     *log.Logger
	prompt  *gprompt.Prompt
	world   *model.World
	history *os.File
}

// Create returns a fresh Shell.
func Create(world *model.World) (*Shell, error) {
	var (
		err     error
		history []string
		buf     bytes.Buffer
		sh      = &Shell{
			world: world,
		}
	)

	if sh.log, err = common.GetLogger(logdomain.Shell); err != nil {
		return nil, err
	} else if sh.history, err = os.OpenFile(
		common.HistPath,
		os.O_RDWR|os.O_CREATE|os.O_APPEND|os.O_SYNC,
		0644); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			history = make([]string, 0)
		} else {
			sh.log.Printf("[ERROR] Cannot open history file %s: %s\n",
				common.HistPath,
				err.Error())
			return nil, err
		}
	} else if _, err = io.Copy(&buf, sh.history); err != nil {
		sh.log.Printf("[ERROR] Failed to read history file %s: %s\n",
			common.HistPath,
			err.Error())
		return nil, err
	} else {
		history = strings.Split(buf.String(), "\n")
	}

	sh.prompt = gprompt.New(
		sh.execute,
		sh.completions,
		gprompt.OptionHistory(history),
		gprompt.OptionTitle(common.AppName),
		gprompt.OptionPrefix(prefix),
		gprompt.OptionLivePrefix(sh.livePrefix),
	)

	return sh, nil
} // func Create(world *model.World) (*Shell, error)

func (sh *Shell) Run() {
	sh.prompt.Run()
}

func (sh *Shell) livePrefix() (string, bool) {
	if sh.world == nil {
		return prefix, true
	}

	var newPrefix = fmt.Sprintf("(%s) %s",
		sh.world.CurrentLocation().Name,
		prefix)

	return newPrefix, true
} // func (sh *Shell) livePrefix() (string, bool)

func (sh *Shell) execute(input string) {
	sh.log.Printf("[DEBUG] Execute input %q\n",
		input)
} // func (sh *Shell) execute(input string)

func (sh *Shell) completions(d gprompt.Document) []gprompt.Suggest {
	// sh.log.Printf("[DEBUG] Attempt to complete input %q\n",
	// 	d.Text)
	return []gprompt.Suggest{}
} // func (sh *Shell) completions(d *gprompt.Document) []gprompt.Suggest
