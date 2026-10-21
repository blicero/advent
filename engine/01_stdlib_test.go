// /home/krylon/go/src/github.com/blicero/advent/01_stdlib_test.go
// -*- mode: go; coding: utf-8; -*-
// Created on 20. 10. 2026 by Benjamin Walkenhorst
// (c) 2026 Benjamin Walkenhorst
// Time-stamp: <2026-10-21 16:25:26 krylon>
//
// This file is the test suite for the "standard library" I provide to the
// Lua engine. It could do with some cleaning up, the test scripts write to
// stdout *a lot*, which makes the output of the test harder to read.
//
// Donnerstag, 17. 08. 2017, 20:51
// This file is getting fairly large. I should consider breaking it up into
// one test for the basic Go-Lua-interface, and then one test for each
// module. Yeah, Coulda, Shoulda, Woulda...
//
// Freitag, 01. 09. 2017, 19:36
// So, I am in the hospital, currently, where it is of course impossible
// to access my home network. I don't even have Internet access, unless
// I use my phone as a Wifi hotspot, but the signal quality around here is
// not always very good.
// In other words, I need to do something about my tests. First of all,
// it must not depend on being able to access my home network. If I were
// to share this code with anyone, that would be a prudent idea, anyway.
// Secondly, to accomodate my current lack of any kind of network access,
// it would be nice if the test suite could still pass without network
// access.
// I could manually start some sockets on localhost.
// That is probably the best idea, or at least the best I can come up with
// right now.
// Let's assume, for starters, that ssh is reachable on localhost.

package engine

import (
	"os"
	"path/filepath"
	"testing"
)

var eng *Engine

// Sonntag, 13. 08. 2017, 02:25
// I need to write a couple of test cases for the TCP interface, right?
// nolint: unused
const luaTCPTestScript = `
-- BW 20160825_1936
-- Baby steps
--

-- Freitag, 01. 09. 2017, 19:44
-- I want my tests to pass cleanly even without access to any particular
-- network, so for the time being, I assume that ssh is reachable on
-- localhost. That is probably not an outrageous assumption.
--conn, err = tcp.connect("192.168.0.1", 21)
conn, err = tcp.connect("127.0.0.1", 22)

if err ~= nil then
   print("Error connecting:", err)
else
   print("Yo:", conn)
end

sent, err = conn:send("Hallo")
if err ~= nil then
  print("Error sending data:", err)
else
  print("Sent Message:", sent)
end

reply,err = conn:receive()

if err ~= nil then
  print("Error:", err)
else
  print("Received a reply:", reply)
end

conn:close()

print("Bye from Lua")

return true
`

const luaRegexTestScript = `
-- BW, 20160826_2009
-- Check it out!

pat = "(\\d+)"
re = regexp.compile(pat)

print(re)

result = re:match("Peter1234")
print(result)

result = re:find("Peter666Lustig")
print(result)

result = re:replace("Peter666Lustig", "123")
print(result)

print("Bye from Lua")

return true
`

// TODO - add test cases for the time interface
// nolint: unused
const luaTimeTestScript = `
-- BW, 20160829_2009

t1 = time.now()
dur = duration.seconds(2)

print(dur)
s = tostring(dur)
print(s)

t2 = t1 + dur

print(tostring(t1), tostring(t2))

s = tostring(t1)

print("Converted time value to string")
print(s)

dur = duration.seconds(2)

time.sleep(dur)

t2 = time.now()

delta = t2 - t1

print("Time passed:", tostring(delta))

unixt, err = time.from_unix(3600)
if err ~= nil then
  print("Cannot create time from Unix stamp: ", err)
  return false, err
end

unixs, err = unixt:unix()
if err ~= nil then
  print("Cannot convert time value to Unix stamp: ", err)
  return false, err
elseif unixs ~= 3600 then
  msg = "Conversion to Unix stamp yielded incorrent value: " .. tostring(unixs)
  print(msg)
  return false, msg
end

return true
`

func TestCreateEngine(t *testing.T) {
	var err error

	if eng, err = NewEngine(); err != nil {
		t.Errorf("Failed to create script engine: %s",
			err.Error())
	}
} // func TestCreateEngine(t *testing.T)

func TestRegexBinding(t *testing.T) {
	if eng == nil {
		t.SkipNow()
	}
	var err error
	eng.Script = luaRegexTestScript

	eng.State.SetTop(0)

	if err = eng.Run(); err != nil {
		t.Error(err.Error())
	}
} // func TestRegexBinding(t *testing.T)

// nolint: gocyclo
func TestSetGlobal(t *testing.T) {
	if eng == nil {
		t.SkipNow()
	}

	// Donnerstag, 10. 08. 2017, 20:24
	// I should add one more test for (light)userdata?
	// But then I need a way to conveniently set
	// meta tables, right?

	eng.State.SetTop(0)

	var err error
	var intVal = 42
	var int8Val int8 = 65
	var int16Val int16 = 16384
	var int32Val int32 = 1048576
	var stringVal = "Wer das liest, ist doof."
	var boolVal = true
	var f32Val float32 = 3
	var f64Val = 2.71828182846

	if err = eng.SetGlobalValue("zahl", intVal, ""); err != nil {
		t.Fatalf("Error setting global value from int: %s",
			err.Error())
	} else if err = eng.SetGlobalValue("text", stringVal, ""); err != nil {
		t.Fatalf("Error setting global value from string: %s",
			err.Error())
	} else if err = eng.SetGlobalValue("wahrheit", boolVal, ""); err != nil {
		t.Fatalf("Error setting global value from bool: %s",
			err.Error())
	} else if err = eng.SetGlobalValue("pi", f32Val, ""); err != nil {
		t.Fatalf("Error setting global value from float32: %s",
			err.Error())
	} else if err = eng.SetGlobalValue("euler", f64Val, ""); err != nil {
		t.Fatalf("Error setting global value from float64: %s",
			err.Error())
	} else if err = eng.SetGlobalValue("i8", int8Val, ""); err != nil {
		t.Fatalf("Error setting global value from int8: %s",
			err.Error())
	} else if err = eng.SetGlobalValue("i16", int16Val, ""); err != nil {
		t.Fatalf("Error setting global value from int16: %s",
			err.Error())
	} else if err = eng.SetGlobalValue("i32", int32Val, ""); err != nil {
		t.Fatalf("Error setting global value from int32: %s",
			err.Error())
	}

	const globalsTestScript = `
if type(zahl) ~= "number" or zahl ~= 42 then 
    return false, "int"
elseif type(text) ~= "string" or text ~= "Wer das liest, ist doof." then
    return false, "string"
elseif type(wahrheit) ~= "boolean" or not wahrheit then
    return false, "bool"
elseif type(pi) ~= "number" or pi ~= 3 then
    return false, "float32"
elseif type(euler) ~= "number" or euler ~= 2.71828182846 then
    return false, "float64"
else
    return true, "nothing"
end

return true
`

	eng.Script = globalsTestScript

	if err = eng.Run(); err != nil {
		t.Error(err.Error())
	}

	// Now we fetch the return values from the script from Lua's stack
	// and check if everything went as planned:
	var res, ok bool
	var msg string

	// Maybe I got the order of the return values wrong?
	// From the Lua reference manual section on the function lua_call(3),
	// I think the leftmost argument is pushed first, so after the script
	// runs, the last return value is on top of the stack.
	if !eng.State.IsString(-1) {
		t.Fatalf("Top value on stack after test script is not a string, but a %s",
			eng.State.TypeOf(-1).String())
	} else if !eng.State.IsBoolean(-2) {
		t.Fatalf("Second value on stack is not a boolean, but a %s",
			eng.State.TypeOf(-2).String())
	}

	if msg, ok = eng.State.ToString(-1); !ok {
		t.Fatal("CANTHAPPEN: Failed to fetch string from the Lua stack")
	}
	res = eng.State.ToBoolean(-2)

	if !res {
		t.Fatalf("At least one global did not arrive unharmed: %s",
			msg)
	}
} // func TestSetGlobal(t *testing.T)

func TestLuaScriptReturnValues(t *testing.T) {
	if eng == nil {
		t.SkipNow()
	}

	var err error

	// In order to test the behavior of eng.Run when dealing
	// with the return values from Lua scripts.
	// If this thing ever sees any significant use, this
	// will probably happen a lot.
	// Plus, I think, it would make more sense to have
	// more flexible return values. A Lua Probe should return
	// the same triplet of values returned by builtin probes.

	eng.State.SetTop(0)
	eng.Script = "return true"
	if err = eng.Run(); err != nil {
		t.Errorf("eng.Run returned an unexpected error (should have return nil: %s",
			err.Error())
	}

	eng.State.SetTop(0)
	eng.Script = "return false"
	if err = eng.Run(); err == nil {
		t.Error("Expected an error value from eng.Run")
	} else if err.Error() != "lua script returned false, but no error message" {
		// This return value is a hard coded constant in the source code.
		t.Errorf("Unexpected error message from LuaEngine: %s",
			err.Error())
	}

} // func TestLuaReturnFalse(t *testing.T)

func TestLuaLogModule(t *testing.T) {
	if eng == nil {
		t.SkipNow()
	}
	var err error
	const luaLogScript = `
if log == nil then
    return false, "Log module is nil"
end

for k, v in pairs(log) do
  print("==> " .. k)
end

return true

`

	eng.State.SetTop(0)

	eng.Script = luaLogScript
	if err = eng.Run(); err != nil {
		t.Error(err.Error())
	}
} // func TestLuaLogModule(t *testing.T)

func TestLuaIncludeFile(t *testing.T) {
	if eng == nil {
		t.SkipNow()
	}

	var err error
	eng.State.SetTop(0)

	eng.Script = `
return include("tst01.lua")
`

	// So far, we have not set an include path, so this call should fail
	if err = eng.Run(); err == nil {
		t.Error("Include without setting an inlude path should fail")
	} else {
		eng.State.SetTop(0)
	}

	var cwd string

	if cwd, err = os.Getwd(); err != nil {
		t.Fatalf("Error getting working directory: %s",
			err.Error())
	}

	eng.LookupFolders = []string{
		filepath.Join(cwd, "testdata"),
	}

	if err = eng.Run(); err != nil {
		t.Fatalf("Error including tst01.lua: %s",
			err.Error())
	}

	var num int
	var str string
	var ok bool

	eng.State.Global("testvar01")
	if num, ok = eng.State.ToInteger(-1); !ok {
		t.Errorf("Global testvar01 should be an integer, not a %s",
			eng.State.TypeOf(-1).String())
	} else if num != 100 {
		t.Errorf("Global testvar01 should equal 100, not %d",
			num)
	}

	eng.State.Pop(1)

	eng.State.Global("testvar02")
	if str, ok = eng.State.ToString(-1); !ok {
		t.Errorf("Global testvar02 should be a string, not a %s",
			eng.State.TypeOf(-1).String())
	} else if str != "Hallo, Zebu!" {
		t.Errorf("Global testvar02 should be \"Hallo, Zebu!\", not %q",
			str)
	}

	eng.State.SetTop(0)

} // func TestLuaIncludeFile(t *testing.T)

func TestIncludeDir(t *testing.T) {
	if eng == nil {
		t.SkipNow()
	}

	var err error

	eng.Script = `
return includedir("testvars")
`

	eng.State.SetTop(0)

	if err = eng.Run(); err != nil {
		t.Errorf("Error running test script for includedir - %s",
			err.Error())
	}

	var num int
	var value, ok bool
	eng.State.Global("testvar03")
	if num, ok = eng.State.ToInteger(-1); !ok {
		t.Errorf("Expected testvar03 to an integer, not a %s",
			eng.State.TypeOf(-1).String())
	} else if num != 15 {
		t.Errorf("Expected testvar03 to equal 15, not %d",
			num)
	}
	eng.State.SetTop(0)

	eng.State.Global("testvar04")
	if !eng.State.IsBoolean(-1) {
		t.Errorf("Global testvar04 should be a boolean value, not a %s",
			eng.State.TypeOf(-1).String())
	} else if value = eng.State.ToBoolean(-1); !value {
		t.Error("Global testvar04 should equal true")
	}
} // func TestIncludeDir(t *testing.T)
