// /home/krylon/go/src/github.com/blicero/advent/logdomain/logdomain.go
// -*- mode: go; coding: utf-8; -*-
// Created on 18. 10. 2026 by Benjamin Walkenhorst
// (c) 2026 Benjamin Walkenhorst
// Time-stamp: <2026-10-18 11:09:37 krylon>

// Package logdomain provides constants for logging.
package logdomain

//go:generate stringer -type=ID

// ID represents a part of the program that wants to write to the log.
type ID uint8

const (
	Parser ID = iota
)

// All returns a slice of all defined ID values.
func All() []ID {
	return []ID{
		Parser,
	}
}
