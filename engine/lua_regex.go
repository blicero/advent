// /home/krylon/go/src/github.com/blicero/advent/engine/lua_regex.go
// -*- mode: go; coding: utf-8; -*-
// Created on 21. 10. 2026 by Benjamin Walkenhorst
// (c) 2026 Benjamin Walkenhorst
// Time-stamp: <2026-10-21 15:51:19 krylon>

package engine

import (
	"fmt"
	"regexp"

	lua "github.com/Shopify/go-lua"
)

func regexCompile(state *lua.State) int {
	var pattern, msg string
	var re *regexp.Regexp
	var err error

	pattern = lua.CheckString(state, 1)

	if re, err = regexp.Compile(pattern); err != nil {
		msg = fmt.Sprintf("Error compiling regex: %s", err.Error())
		fmt.Println(msg)
		state.PushNil()
		state.PushString(msg)
		return 2
	}

	state.PushUserData(re)
	lua.SetMetaTableNamed(state, "re_meta")
	return 1
} // func regex_compile(state *lua.State) int

func regexFind(state *lua.State) int {
	var lre interface{}
	var re *regexp.Regexp
	var instance string

	lre = lua.CheckUserData(state, 1, "re_meta")
	instance = lua.CheckString(state, 2)

	switch v := lre.(type) {
	case *regexp.Regexp:
		re = v
	default:
		panic("Unexpected type for regex!")
	}

	result := re.FindStringSubmatch(instance)

	if result == nil {
		return 0
	}

	var cnt = len(result)

	for _, sub := range result {
		state.PushString(sub)
	}

	return cnt
} // func regex_find(state *lua.State) int

func regexMatch(state *lua.State) int {
	var lre interface{}
	var re *regexp.Regexp
	var instance string

	lre = lua.CheckUserData(state, 1, "re_meta")
	instance = lua.CheckString(state, 2)

	switch v := lre.(type) {
	case *regexp.Regexp:
		re = v
	default:
		panic("I EXPECTED A REGEXP!")
	}

	state.PushBoolean(re.MatchString(instance))
	return 1
} // func regex_match(state *lua.State) int

func regexReplace(state *lua.State) int {
	var lre interface{}
	var re *regexp.Regexp
	var orig, replacement, changed string

	lre = lua.CheckUserData(state, 1, "re_meta")
	orig = lua.CheckString(state, 2)
	replacement = lua.CheckString(state, 3)

	switch v := lre.(type) {
	case *regexp.Regexp:
		re = v
	default:
		panic("Userdata is not a regexp!")
	}

	changed = re.ReplaceAllString(orig, replacement)

	state.PushString(changed)

	return 1
} // func regex_replace(state *lua.State) int
