#!/usr/bin/env lua
-- Time-stamp: <2026-10-22 07:48:05 krylon>
--
-- /media/disk/lua/util.lua
-- created on 06. 09. 2008
-- (c) 2008 Benjamin Walkenhorst

module("krylib", package.seeall)

--- A collection of general-purpose functions
--- @author Benjamin Walkenhorst (krylon (at) gmx (dot) net)

-- require("posix")

--- Return true if key is present in tab
--- @param tab A table
--- @param key The key to look for
--- @return true if the given key is present in tab
function has_key(tab, key)
   assert(is_table(tab))
   return (nil ~= tab[key])
end -- function has_key(tab, key)

--- Return true if x is not nil
--- @param x A Lua value
--- @return true if x is not nil 
function defined(x)
   return nil ~= x
end

--- Return true if x is nil
--- @param x A Lua value
--- @return True if x is nil
function is_null(x)
   return nil == x
end

--- Return true if n is greater than zero.
--- @param n A Number
--- @return True if n > 0
function is_positive(n)
   return n > 0
end

--- Return true if n is zero
--- @param n A Number
--- @return True if n == 0
function is_zero(n)
   return n == 0
end

--- Return true if n is less than zero
--- @param n A Number
--- @return True if n < 0
function is_negative(n)
   return n < 0
end

--- Check if x is a function
--- @param x A Lua value
--- @return true if x is a function
function is_function(x)
   return type(x) == "function"
end

--- Check if x is a thread
--- @param x A Lua value
--- @return true if x is a thread 
function is_thread(x)
   return type(x) == "thread"
end

--- Check if x is a table
--- @param x A Lua value
--- @return true if x is a table
function is_table(x)
   return type(x) == "table"
end

--- Check if x is a string
--- @param x A Lua value
--- @return true if x is a string
function is_string(x)
   return type(x) == "string"
end

--- Convert a boolean value to a string.
--- @param x A boolean value
--- @return The string representation of the boolean value,
---         or nil if the parameter is not a boolean.
function bool_to_str(token)
   if true == token then
      return "true"
   elseif false == token then
      return "false"
   else
      return nil
   end
end

function str_to_bool(token)
   if token == "true" then
      return true
   elseif token == "false" then
      return false
   else
      return nil
   end
end

--- Check if the parameter is a number.
--- @param x A Lua value
--- @return true if the parameter is a number. 
function is_number(item)
   return type(item) == "number"
end -- function is_number(item)

--- Return the absolute value of the difference between two numbers.
function absdiff(a, b)
   assert(is_number(a))
   assert(is_number(b))

   return math.abs(a - b)
end -- function absdiff(a, b)

--- Return an iterator over a table's values.
--- I don't think I ever tested this or that it works
--- at all.
--- This is meant to be used on arrays.
function each(tab)
   local i = 0
   local n = table.getn(tab)
   return function()
             i = i + 1
             if i <= n then return tab[i] end
          end
end -- function each(n)

--- Return an iterator over the values of a table.
--- Unlike each(), this works for all kinds of tables.
function values(tbl)
   local ikey = nil
   return function()
	     local nkey, nval
	     nkey, nval = next(tbl, ikey)
	     if nval then
		ikey = nkey
		return nval
	     end
	  end
end -- function values(tbl)

--- Return an array of all keys in the given table.
--- @param tbl A Table
--- @return An array of all keys in the given table. 
function keys(tbl)
   local k = {}
   for key, _ in pairs(tbl) do
      table.insert(k, key)
   end
   return k
end -- function keys(tbl)

--- Return a table of the values between a and b.
--- @param a A number
--- @param b A number
--- @return A Table of integers in the interval [a, b] 
function range(a, b)
   local r = {}
   assert(is_number(a))
   assert(is_number(b))
   assert(a <= b)

   for i = a, b, 1 do
      table.insert(r, i)
   end

   return r
end

--- Return a coroutine that iterates over the integers between a and b.
--- So redundant!
-- function irange(a, b)
--    return coroutine.create(function ()
--                               local i = a
--                               while i <= b do
--                                  coroutine.yield(i)
--                                  i = i + 1
--                               end
--                            end)
-- end

--- Check if item is of the expected type.
--- So redundant!
function typecheck(item, expected)
   return type(item) == expected
end

--- Check if item is of the expected type, raise an error if it is not.
function etypecheck(item, expected)
   return assert(type(item) == expected,
                 string.format("Expected type %s, but item is a %s",
                               expected, type(item)))
end

--- Split a string on occurrences of pattern.
--- IMCOMPLETE!
function split(str, pat)
   local a, b, parts = 0, 0, {}
   a = string.find(str, pat)
   if is_number(a) then
   end
end

--- Create a closure that returns monotonically increasing numbers.
function make_counter_fn(start, step)
   local counter = start
   local function counter_function()
      counter = counter + step
      return counter
   end
   return counter_function
end

--[[
Sonntag, 14. 06. 2009, 15:50
Es wäre schön, hier auch den Namen der Funktion zu kennen.
Wenn ich so etwas wie ... Den Namen und die Tabelle übergebe -
selbst globale Variablen sind ja eigentlich nur Key in der
Tabelle _G - dann könnte ich schauen, ob ... dingsen.

Sonntag, 14. 06. 2009, 15:57
Hihi, es geht. Ich müsste nur zusehen, dass ich irgendwie mit
mehrfachen Rückgabewerten eleganter umgehe, und es wäre gut, die
Funktion auch wieder zur Laufzeit untracen zu können.
Dafür sollte ich das wahrscheinlich am besten in ein eigenes
Modul auslagern.
Streng genommen schaue ich mir gerade noch einmal ein Profiler-Modul an, das
ich mal geschrieben habe, und es es sieht so aus, als könnte ich das
Debugging-Interface von Lua dafür verwenden, was vermutlich auch
cleverer ist.
--]]

--- Trace the function stored under the given name in the given table.
function trace_function(function_name, table)
   if is_null(table) then table = _G end
   -- Da ich den String jedes Mal mit ausgebe, ist es
   -- wohl besser, den einmal zu erzeugen und zu recyclen.
   local trace_str = "Calling " .. function_name .. ":"
   local orig = table[function_name]
   
   local function traced_fn(...)
      print(trace_str)
      print("\tArgs:", ...)
      local result = orig(...)
      print("\tResult:", result)
      return result
   end
   
   table[function_name] = traced_fn
   return orig
end

-- function make_array(...)
--    return {...}
-- end

-- Samstag, 06. 03. 2010, 02:19
-- Ich muss mir noch mal was überlegen, wie ich geschachtelte Tabellen
-- ansehnlich ausgeben kann.

--- Print the contents of a table.
function dump_table(tbl)
   for k,v in pairs(tbl) do
      print(k, "=>", v)
   end
end

--- Return the sign of a number.
function sign(n)
   assert(is_number(n))
   if n < 0 then
      return -1
   elseif n == 0 then
      return 0
   elseif n > 0 then
      return 1
   end
end

--- Return a merged table containing union of the keys in t1 and t2.
--- If a key is present in both tables, the returned table will contain
--- the value from t2.
--- Shallow copy!
function merge_tables(t1, t2)
   assert(is_table(t1))
   assert(is_table(t2))
   local combined = {}
   for key, value in ipairs(t1) do
      combined[key] = value
   end
   for key, value in ipairs(t2) do
      combined[key] = value
   end
   return combined
end

--- Format a number of bytes in a human readable format,
--- i.e. 8192 becomes 8.000 KiB,
--- 1048576 becomes 1.000 MiB
function fmt_bytes(num_bytes)
   assert(is_number(num_bytes))
   local units = { "Bytes", "KiB", "MiB", "GiB", "TiB", "PiB", }
   local i = 1
   while num_bytes >= 1023 and i <= table.maxn(units) do
      num_bytes = num_bytes / 1024
      i = i + 1
   end
   return string.format("%.3f %s", num_bytes, units[i])
end

--[[
The Art Of UNIX Programming tells us how to turn a process into a
proper daemon. There are 6 steps involved:
1 set the umask to 0
2 fork
3 change to root directory
4 call setsid(2) to start a new session
- possibly fork again to avoid being a session leader
(on SysV, a session leader might be re-attached to a
 controlling terminal)
5 close unnecessary file descriptors
6 close STD* file descriptors
Since we know nothing about the file descriptors of the calling process,
we leave that to the caller... Unless silence=False is passed to the
function, STD(OUT|ERR|IN) will be reopened on /dev/null.
   --]]

--- Daemonize the current process.
--- NOT FINISHED!
-- function daemonize()
--    -- Freitag, 10. 06. 2011, 21:44
--    -- Ich hoffe, ich habe die Syntax richtig hin bekommen...
--    posix.chmod("u=rwx g=rx a=rx") -- Step 1

--    local pid = posix.fork()
--    if pid == -1 then -- Something went horribly wrong!
--       error(string.format("Error forking process: %s", posix.errno()))
--    elseif pid == 0 then -- Child process
--       -- do nothing
--    else -- Parent process
--       os.exit(0)
--    end

--    -- Dabei kann doch nun wirklich nichts schief gehen, oder?
--    posix.chdir("/")

--    -- Lua has no setsid() bindings

--    -- How do we know about unnecessary file descriptors?

--    -- Close STD*
-- end -- function daemonize()

--- Expand ~ in the given string to user's home directory, i.e. if the current
--- user's home directory was /home/krylon, then ~/.emacs.d/init.el becomes
--- /home/krylon/.emacs.d/init.el
function expanduser(path)
   assert(typecheck(path, "string"))
   local abs_path, cnt = string.gsub(path, "~", os.getenv("HOME"), 1)
   return abs_path
end -- function expanduser(path)

--- Return the name of the operating system we're running on.
--- This calls uname and reads its output, so this only works on systems that
--- make a reasonable effort to behave like a Unix system.
function osname()
   local proc = io.popen("/bin/uname")
   local os_name = proc:read()
   proc:close()
   return os_name
end -- function osname()

--- Roll print and string.format into one function, like printf() from
--- the C standard library.
function printf(...)
   print(string.format(...))
end -- function printf(...)

function dbg_printf(fmt, ...)
   local status, msg = pcall(string.format, fmt, ...)
   if status then
      print(msg)
   else
      print(fmt, ...)
      print(debug.traceback(string.format("Error formating string: %s", msg)))
      debug.debug()
   end
end -- function dbg_printf(fmt, ...)

--- Print the formatted message to stdout, prefixing it with the current time
--- and date.
function log_msg(fmt, ...)
   local stamp = os.date()
   dbg_printf("[%s] " .. fmt, stamp, ...)
end -- function log_msg(fmt, ...)

--- Remove a trailing line break from a given string
--- @param s A string
--- @return The string with a trailing line break ("\r?\n") removed.
function strip(s)
   return string.gsub(s, "\r?n$", "")
end -- function strip(s)

--- Return the hostname of the system we're running on.
--- This functions works by using io.popen to call hostname,
--- so it only works on systems that make a reasonable effort to behave like
--- Unix systems.
function gethostname()
   local proc = io.popen("/bin/hostname")
   local hostname = proc:read()
   proc:close()
   return hostname
end -- function gethostname()

--- Return a UUID
--- This function works by using io.popen to call uuidgen, which I found both on
--- my Ubuntu boxes as well as on my FreeBSD-VM and my NetBSD system, so I naively
--- assume it to be quasi-standard. It is not mandated by POSIX, however...
--- @return A UUID string.
function get_uuid()
   local proc = io.popen("uuidgen")
   local uuid = proc:read()
   proc:close()
   return uuid
end -- function get_uuid()

function curry1(fn, a)
   return function(...) fn(a, ...) end
end -- function curry1(fn, a)

function get_zone_for_host(name)
   assert(type(name) == "string", "name must be a string!")
   return string.match(name, "^[^.]+[.](.*)")
end -- function get_zone_for_host(hostname)

function is_tld(zone)
   return (type(zone) == "string") and string.match(zone, "^[^.]+[.]?$")
end -- function is_tld(zone)

function IMPLEMENT_ME()
   error("IMPLEMENT ME!", 2)
end -- function IMPLEMENT_ME()

--- Return the average of the numbers in the given array.
--- This function does not check if the values in the array
--- actually are numbers, it is the reponsibility of the caller
--- to make sure to pass in only valid arrays.
--- @param tbl An array of numbers.
--- @return The average of all the numbers in the array.
function average(tbl)
   local sum = 0
   local cnt = 0
   for _, value in ipairs(tbl) do
      sum = sum + value
      cnt = cnt + 1
   end -- for _, value in ipairs(tbl) do
   return sum / cnt
end -- function average(tbl)

-- function is_file(path)
--    local status, msg = posix.access(path, "f")
--    if nil == status then
--       return false
--    elseif 0 == status then
--       return true
--    else
--       printf("Unexpected return value from posix.access: %s", tostring(status))
--       return false
--    end
-- end -- function is_file(path)

-- function isdir(path)
--    local status, msg = posix.stat(path, "type")
--    if nil == status then
--       return nil, string.format("Could not stat %s: %s", path, msg)
--    else
--       return status == "directory"
--    end
-- end -- function isdir(path)

function path_join(root, suffix)
   local SEP = "/"
   if string.find(root, "/$") then
      return root .. suffix
   else
      return root .. SEP .. suffix
   end
end -- function path_join(root, suffix)

function irange(a, b, step)
   local i = a
   local iter = function()
                   if i <= b then
                      local r = i
                      i = i + step
                      return r
                   else
                      return nil
                   end
                end

   return iter
end -- function irange(a, b, step)

function looks_like_ip(s)
   assert(type(s) == "string")

   b1, b2, b3, b4 = string.match(s, "^(%d+)[.](%d+)[.](%d+)[.](%d+)$")

   if b1 == nil then
      return false
   else
      return (b1 >= 0 and b1 < 256 and
              b2 >= 0 and b2 < 256 and
              b3 >= 0 and b3 < 256 and
              b4 >= 0 and b4 < 256)
   end -- if b1 == nil then
end -- function looks_like_ip(s)
