// /home/krylon/go/src/github.com/blicero/advent/model/model.go
// -*- mode: go; coding: utf-8; -*-
// Created on 18. 10. 2026 by Benjamin Walkenhorst
// (c) 2026 Benjamin Walkenhorst
// Time-stamp: <2026-10-19 09:07:27 krylon>

// Package model provides data types used throughout the program.
package model

// Location is a place in the game.
type Location struct {
	ID          int64
	Name        string
	Description string
}

// Character is a sentient being the player can communicate and interact with.
type Character struct {
	ID          int64
	Name        string
	Description string
}

// Object is an inanimate thing, what physisicts call a "body".
type Object struct {
	ID          int64
	Name        string
	Description string
}

// World is the entirety of the Game.
type World struct {
	Places []Location
	People []Character
	Things []Object
	curLoc int64
}

func (w *World) CurrentLocation() Location {
	return w.Places[w.curLoc]
}
