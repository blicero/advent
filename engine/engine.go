// /home/krylon/go/src/github.com/blicero/advent/engine/engine.go
// -*- mode: go; coding: utf-8; -*-
// Created on 20. 10. 2026 by Benjamin Walkenhorst
// (c) 2026 Benjamin Walkenhorst
// Time-stamp: <2026-10-20 18:12:48 krylon>

// Package engine provides bindings to the Lua engine.
package engine

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	lua "github.com/Shopify/go-lua"
	"github.com/blicero/advent/common"
	"github.com/blicero/advent/logdomain"
	"github.com/blicero/krylib"
	"github.com/davecgh/go-spew/spew"
)

// Engine wraps the Lua VM and associated state.
type Engine struct {
	log           *log.Logger
	luaLog        *log.Logger // nolint: megacheck,structcheck
	State         *lua.State
	LookupFolders []string
	Script        string
}

// NewEngine returns a fresh ScriptEngine using the given folders to look
// for script files.
func NewEngine(includeDir ...string) (*Engine, error) {
	var (
		err error
		eng = &Engine{
			State:         lua.NewState(),
			LookupFolders: includeDir,
		}
	)

	if eng.log, err = common.GetLogger(logdomain.Lua); err != nil {
		return nil, err
	}

	return eng, nil
} // func NewEngine(includeDir ...string) (*Engine, error)

func (eng *Engine) initialize() error {
	lua.OpenLibraries(eng.State)
	eng.State.PushBoolean(common.Debug)
	eng.State.SetGlobal("DEBUG")

	return nil
} // func (eng *ScriptEngine) initialize() error

// AddFunction registers a function (written in Go) with the Lua VM,
// globally, under the given name.
// If there is a previous binding for the name, it is overwritten.
// Use with care!
func (eng *Engine) AddFunction(name string, fn func(*lua.State) int) error {
	eng.State.Register(name, fn)
	return nil
} // func (eng *Engine) AddFunction(name string, fn func(*lua.State) int) error

// PushGoValue pushes one value from Go on top of the Lua stack.
// This method uses reflection to figure out the type of the argument.
// Anything that cannot be converted straight to a Lua type is pushed
// as a LightUserData.
// If meta is a non-empty string, it is taken as the name of the metatable
// to associate with the value.
// nolint: gocyclo
func (eng *Engine) PushGoValue(val interface{}, meta string) error {
	if val == nil {
		// This is a special case, insofar as "nil" can be nil of many
		// types in Go, in Lua, it does not matter, or rather nil is its
		// own type
		eng.State.PushNil()
		return nil
	}

	// Next on the list would be Table, but I am not sure how I could
	// convert that without major contortions.
	//
	// XXX I think I need to handle unsigned integer values separately,
	//     at least uint64, because I am not sure what happens when I
	//     convert that to an int.
	//     --> There is a separate method "PushUnsigned" that takes care of
	//         the unsigned values.
	switch v := val.(type) {
	case bool:
		eng.State.PushBoolean(v)
	case int:
		eng.State.PushInteger(v)
	case int8:
		eng.State.PushInteger(int(v))
	case int16:
		eng.State.PushInteger(int(v))
	case int32:
		eng.State.PushInteger(int(v))
	case int64:
		eng.State.PushInteger(int(v))
	case uint:
		eng.State.PushUnsigned(v)
	case uint8:
		eng.State.PushUnsigned(uint(v))
	case uint16:
		eng.State.PushUnsigned(uint(v))
	case uint32:
		eng.State.PushUnsigned(uint(v))
	case uint64:
		eng.State.PushUnsigned(uint(v))
	case float32:
		eng.State.PushNumber(float64(v))
	case float64:
		eng.State.PushNumber(v)
	case string:
		eng.State.PushString(v)
	default:
		eng.State.PushLightUserData(val)
		// var msg = fmt.Sprintf("I do not know how to push type %T on the Lua stack",
		// 	val)
		// eng.log.Println("[ERROR] " + msg)
		// return errors.New(msg)
	}

	if meta != "" {
		// Maybe I should check if the given name actually refers to a
		// metatable known to the Lua engine?
		lua.SetMetaTableNamed(eng.State, meta)
	}

	return nil
} // func (eng *Engine) PushGoValue(val interface{}, meta string) error

// SetGlobalValue pushes a value on the Lua stack and, if successful,
// registers it as a global under the given name.
func (eng *Engine) SetGlobalValue(name string, val interface{}, meta string) error {
	if err := eng.PushGoValue(val, meta); err != nil {
		return err
	}

	eng.State.SetGlobal(name)
	return nil
} // func (eng *Engine) SetGlobalValue(name string, val interface{}) error

// Run executes the the engine's Script.
func (eng *Engine) Run() error {
	var err = lua.DoString(eng.State, eng.Script)
	var msg string

	if err != nil {
		eng.log.Printf("[ERROR] Execution of Lua code ended with error code: %s",
			err.Error())
		return err
	}

	// I need to check the return value(s) from the script!
	var retVal, ok bool
	var errmsg string
	var top = eng.State.Top()

	// Does it even make sense to check that the number of values on the
	// stack is *exactly* 1 or 2?
	if top < 1 || top > 2 {
		eng.Dumpstack()
		msg = fmt.Sprintf("Unexpected number of arguments: %d (expected either 1 or 2)",
			top)
		// eng.State.PushString(msg)
		// eng.State.Error()
		eng.log.Println("[ERROR]" + msg)
		return errors.New(msg)
	} else if !eng.State.IsBoolean(1) {
		eng.Dumpstack()
		msg = fmt.Sprintf("First return value from Lua script is %s, not bool\n",
			eng.State.TypeOf(1).String())
		eng.log.Println("[ERROR] " + msg)
		return errors.New(msg)
	} else if retVal = eng.State.ToBoolean(1); !retVal && top < 2 {
		eng.Dumpstack()
		msg = "Lua script returned false, but no error message"
		eng.log.Println("[ERROR] " + msg)
		return errors.New(msg)
	} else if retVal {
		eng.Dumpstack()
		return nil
	} else if errmsg, ok = eng.State.ToString(2); !ok {
		eng.Dumpstack()
		msg = fmt.Sprintf("Second return value from Lua script is a %s, we expected a string\n",
			eng.State.TypeOf(2).String())
		eng.log.Println("[ERROR] " + msg)
		return errors.New(msg)
	}

	eng.log.Println("[ERROR] " + errmsg)
	return errors.New(errmsg)
} // func (eng *Engine) Run() error

// nolint: gocyclo
func (eng *Engine) luaIncludeDir(state *lua.State) int {
	var err error
	var path, msg string
	var ok bool

	if common.Debug {
		defer eng.Dumpstack()
	}

	if state.Top() != 1 {
		state.PushString("Invalid number of arguments: Expected 1(!)")
		state.Error()
	} else if path, ok = state.ToString(1); !ok {
		msg = fmt.Sprintf("Argument must be a string, not a %s",
			state.TypeOf(1).String())
		eng.log.Println("[ERROR] " + msg)
		state.PushString(msg)
		state.Error()
	}

	state.SetTop(0)

	// If the argument is a relative path, we try to look it up relative to our
	// list of directories:
	if !filepath.IsAbs(path) {
		for _, dir := range eng.LookupFolders {
			rpath := filepath.Join(dir, path)
			if ok, err = krylib.IsDir(rpath); err != nil {
				msg = fmt.Sprintf("Error checking if %s is a directory: %s",
					rpath,
					err.Error())
				eng.log.Println("[ERROR] " + msg)
			} else if ok {
				path = rpath
				goto PROCESS_FOLDER
			}
		}

		// If we arrive here, the folder to be included could not be found.
		state.PushBoolean(false)
		state.PushFString("Folder %s was not found", path)
		return 2
	} else if ok, err = krylib.IsDir(path); err != nil {
		msg = fmt.Sprintf("Error checking if %s is a directory: %s",
			path,
			err.Error())
		eng.log.Println("[ERROR] " + msg)
		state.PushString(msg)
		state.Error()
	} else if !ok {
		state.PushBoolean(false)
		state.PushFString("Directory %s does not exist",
			path)
		return 2
	}

PROCESS_FOLDER:
	var dirh *os.File
	var files []os.FileInfo

	if dirh, err = os.Open(path); err != nil {
		msg = fmt.Sprintf("Cannot open directory %s: %s",
			path,
			err.Error())
		eng.log.Println("[ERROR] " + msg)
		state.PushString(msg)
		state.Error()
	}

	defer dirh.Close() // nolint: errcheck,gas

	if common.Debug {
		eng.log.Printf("[DEBUG] Include directory %s\n",
			path)
	}

	if files, err = dirh.Readdir(-1); err != nil {
		msg = fmt.Sprintf("Cannot read directory %s: %s",
			path,
			err.Error())
		eng.log.Println("[ERROR] " + msg)
		state.PushString(msg)
		state.Error()
	}

	for _, finfo := range files {

		// Currently, we do not handle recursion.
		if finfo.IsDir() {
			continue
		}

		var fullPath = filepath.Join(path, finfo.Name())

		// if eng.isBlacklisted(fullPath) {
		// 	continue
		// }

		if common.Debug {
			eng.log.Printf("[DEBUG] Process file %s\n",
				fullPath)
		}

		eng.State.SetTop(0)

		// Do I want to or need to filter by some criterion?
		// Not for now.
		if err = lua.DoFile(eng.State, fullPath); err != nil {
			msg = fmt.Sprintf("Error processing %s: %s",
				fullPath,
				err.Error())
			eng.log.Println("{ERROR] " + msg)
			state.PushString(msg)
			state.Error()
		}
	}

	state.PushBoolean(true)
	return 1
} // func (eng *Engine) luaIncludeDir(state *lua.State) int

func (eng *Engine) luaIncludeFile(state *lua.State) int {
	var err error
	var path, msg string
	var ok bool

	if common.Debug {
		defer eng.Dumpstack()
	}

	if eng.State.Top() != 1 {
		msg = fmt.Sprintf("Invalid number of arguments (expected 1): %d",
			eng.State.Top())
		eng.log.Println("[ERROR] " + msg)
		eng.State.PushString(msg)
		eng.State.Error()
	} else if path, ok = eng.State.ToString(1); !ok {
		msg = fmt.Sprintf("Invalid argument type: expected a string, not a %s",
			eng.State.TypeOf(1).String())
		eng.log.Println("[ERROR] " + msg)
		eng.State.PushString(msg)
		eng.State.Error()
	}

	state.SetTop(0)

	// FIXME - I need to iterate over the include directories to check
	//         if the file is in one of those folders.

	if filepath.IsAbs(path) {
		// How do I handle this? Should I possible prohibit this for
		// security reasons?
		eng.State.PushBoolean(false)
		eng.State.PushFString("Invalid file name: %s - absolute paths are not allowed",
			path)
		return 2
	}

	for _, folder := range eng.LookupFolders {
		var fullPath = filepath.Join(folder, path)
		if common.Debug {
			eng.log.Printf("[DEBUG] Trying to include %s\n",
				fullPath)
		}

		if ok, err = krylib.Fexists(fullPath); err != nil {
			msg = fmt.Sprintf("Error checking if %s exists: %s",
				fullPath,
				err.Error())
			eng.log.Println("[ERROR] " + msg)
			continue
		} else if !ok {
			continue
		}

		if err = lua.DoFile(eng.State, fullPath); err != nil {
			msg = fmt.Sprintf("Error executing %s: %s",
				fullPath,
				err.Error())
			eng.log.Println("[ERROR] " + msg)
			state.PushBoolean(false)
			state.PushString(msg)
			return 2
		}

		state.SetTop(0)
		state.PushBoolean(true)
		return 1
	}

	state.SetTop(0)
	state.PushBoolean(false)
	state.PushFString("File %s was not found in any of the include directories",
		path)
	return 2
} // func (eng *Engine) luaIncludeFile(state *lua.State) int

// Input: a string or a table of strings.
// Output: boolean / error?
//
// Do I need to adjust the state's INCLUDE variable? Really, does that thing
// have any purpose at all?
// Not really, I guess, unless I add a metatable and override __newindex
// to basically invoke this function. ???
// However, if I do that, I need to take care handling the table...
// ... ??? I'll have to sleep on that, I guess.
func (eng *Engine) luaIncludeDirAdd(state *lua.State) int {
	var err error
	var msg string
	var path string

	if state.IsString(1) {
		// We already checked if it's a string, so there is no
		// way this can go wrong. None that I can think of,
		// anyway.
		path, _ = state.ToString(1)
		if err = LuaIncludePathAdd(path); err != nil {
			msg = fmt.Sprintf("Error adding path %s to include path: %s",
				path,
				err.Error())
			eng.log.Println("[ERROR] " + msg)
			state.PushBoolean(false)
			state.PushString(msg)
			return 2
		}

		eng.LookupFolders = append(eng.LookupFolders, path)
		state.PushBoolean(true)
		return 1
	} else if state.IsTable(1) {
		// Freitag, 01. 09. 2017, 20:18
		// XXX I am too lazy right now, but I could also check the
		//     length of the table an allocate the slice with the
		//     correct size.
		//     If performance ever becomes an issue, ... yeah.
		var folders = make([]string, 0, 10)
		state.PushNil()
		for state.Next(1) {
			if state.IsString(-1) {
				path, _ = state.ToString(-1)
				folders = append(folders, path)
				state.Pop(1)
			} else {
				// Since we are collecting all the folders
				// before calling LuaIncludePathAdd, I can abort
				// right here, no harm done.
				msg = fmt.Sprintf("Expected string or table of strings. tbl[%v] is a %s",
					state.ToValue(-2),
					state.TypeOf(-1).String())
				eng.log.Println("[ERROR] " + msg)
				state.PushBoolean(false)
				state.PushString(msg)
				return 2
			}
		}

		eng.LookupFolders = append(eng.LookupFolders, folders...)
		state.PushBoolean(true)
		return 1
	} else {
		msg = fmt.Sprintf("Invalid argument type: Expected string or table (of strings), not %s",
			state.TypeOf(1).String())
		eng.log.Println("[ERROR] " + msg)
		state.PushBoolean(false)
		state.PushString(msg)
		return 2
	}
} // func (eng *Engine) luaIncludeDirAdd(state *lua.State) int

// Dumpstack dumps the contents of the Lua stack to standard output, in a
// human-readable format.
func (eng *Engine) Dumpstack() {
	for i := eng.State.Top(); i > 0; i-- {
		// Trailing newline is not required, apparently.
		fmt.Printf("Stackpos #%d: %s",
			i,
			spew.Sdump(eng.State.ToValue(i)))
	}
} // func (eng *Engine) Dumpstack()

// My plan is to tweak it such that Lua call this function like this:
// LogLevel, FormatString, ...FormatArgs
// Is that possible? Can I have variadic functions in Lua?
// ... Seems I can.
// In dem Fall ist LogLevel an Position 1 auf dem Stack (der "untersten" quasi),
// dann FormatString auf Position 2, und alles darüber sind dann die
// Werte, die formatiert werden sollen.
func luaLog(state *lua.State) int {
	var lvl, fmtString, msg string
	var cnt int
	var ok bool

	if cnt = state.Top(); cnt < 2 {
		msg = fmt.Sprintf("Invalid number of arguments: %d (expected >= 2)",
			cnt)
		state.PushString(msg)
		state.Error()
	} else if lvl, ok = state.ToString(1); !ok {
		msg = fmt.Sprintf("First parameter should be a string, not a %s",
			state.TypeOf(1).String())
		state.PushString(msg)
		state.Error()
	} else if !isValidLogLevel(lvl) {
		var levelValues = make([]string, len(common.LogLevels))
		for i, v := range common.LogLevels {
			levelValues[i] = string(v)
		}
		msg = fmt.Sprintf("\"%s\" is not a valid log level. Acceptable values are %s",
			lvl,
			strings.Join(levelValues, ", "))
		state.PushString(msg)
		state.Error()
	} else if fmtString, ok = state.ToString(2); !ok {
		msg = fmt.Sprintf("Second parameter must be a string, not a %s",
			state.TypeOf(2).String())
		state.PushString(msg)
		state.Error()
	}

	// Now, I need to read the remaining arguments into a slice/array
	// and pass that to fmt.Sprintf.
	// Mmmh, that is not going to be entirely trivial, now, is it?

	var args = make([]interface{}, cnt-2)

	for i := 0; i < cnt-2; i++ {
		args[i] = state.ToValue(i + 2)
	}

	var logmsg = fmt.Sprintf(fmtString, args...)

	fmt.Printf("[%s] %s\n",
		lvl,
		logmsg)

	return 0
} // func luaLog(state *lua.Lua) int

///////////////////////////////////////////////////////////////////////
// Utility functions  /////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////

func isValidLogLevel(lvl string) bool {
	// This is a little bit brute-force, but I think the number of
	// log levels is so small, it won't hurt too much.
	for _, name := range common.LogLevels {
		if string(name) == lvl {
			return true
		}
	}

	return false
} // func isValidLogLevel(lvl string) bool

var luaIncludePath = []string{}
var luaIncludeLock sync.RWMutex

// LuaIncludePathAdd adds a list of directories to the luaIncludePath.
// This list of directories will be added to the LookupFolders
// of every new Lua Engine.
func LuaIncludePathAdd(path ...string) error {
	luaIncludeLock.Lock()
	defer luaIncludeLock.Unlock()

	var folders = make(map[string]bool)
	for _, folder := range luaIncludePath {
		folders[folder] = true
	}

	// TODO I probably should check for missing and
	//      duplicate folders.
	//
	//      Now that I think about it - does it even make sense
	//      to check if a folder exists? After all, folders
	//      could be created or deleted after a they have been
	//      added to this list... :-?
	//luaIncludePath = append(luaIncludePath, path...)
	for _, folder := range path {
		var exists bool
		var err error

		if folders[folder] { // Folder is already in include path!
			continue
		} else if exists, err = krylib.IsDir(folder); err != nil {
			// Panicking is a bit overkill, I guess, but I am too
			// lazy right now to think of something better.
			return err
		} else if !exists {
			// Damn it, maybe I should pass in a Logger or something?
			// OR return an error value?
			// Actually, I probably do not want to abort here,
			// just emit some kind of log message?
			continue
		}

		luaIncludePath = append(luaIncludePath, folder)
	}

	return nil
} // func LuaIncludePathAdd(path ...string)

// LuaIncludePathGet returns a copy of the include list.
// NB that the return value is a "deep" copy - changes to it
// will not be reflected in the global include path.
// To make changes to the global include path, please use
// LuaIncludePathAdd and LuaIncludePathDel
func LuaIncludePathGet() []string {
	luaIncludeLock.RLock()

	var folders = make([]string, len(luaIncludePath))
	copy(folders, luaIncludePath)

	luaIncludeLock.RUnlock()
	return folders
} // func LuaIncludePathGet() []string

// LuaIncludePathDel remove the given path from the list of include directories,
// if it is in the list.
// Passing in a path that is already not in the include list is not an error,
// but no special care is taken to check for that case, so the function will
// still take some time to iterate over the include list.
func LuaIncludePathDel(path string) {
	var newList []string

	luaIncludeLock.Lock()
	defer luaIncludeLock.Unlock()

	for idx, folder := range luaIncludePath {
		if path == folder {
			newList = make([]string, len(luaIncludePath)-1)
			// I want to be a little clever about this to avoid
			// making it more inefficient than it needs to be.
			// On the other hand, I smell a off-by-one-error.
			// (At least, in Go, this is easy to test and debug,
			// I remember how much fun I had with these in C...)
			copy(newList[:idx], luaIncludePath[:idx])
			copy(newList[idx:], luaIncludePath[idx+1:])
		}
	}
} // func LuaIncludePathDel(path string)

// LuaIncludePathLookup attempts to find the named file in the lookup directories.
// If successful, it returns the content of that file as a string.
func LuaIncludePathLookup(name string) (string, error) {
	// I am not sure if it's such a good idea to copy
	luaIncludeLock.RLock()
	var folders = make([]string, len(luaIncludePath))
	copy(folders, luaIncludePath)
	luaIncludeLock.RUnlock()

	for _, folder := range folders {
		var fullPath = filepath.Join(folder, name)
		var content []byte
		var err error

		if content, err = os.ReadFile(fullPath); err != nil {
			if err == os.ErrNotExist {
				continue
			} else {
				return "", err
			}
		} else {
			return string(content), nil
		}
	}

	return "", nil
} // func LuaIncludePathLookup(name string) (string, error)
