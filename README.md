gi: a go interpreter
=======
`gi` aims at being a scalable Go REPL
for doing interactive coding and data analysis.
Currently it is backed by LuaJIT, a tracing Just-in-time
compiler for Lua that provides quite
nice performance, sometimes even beating
ahead-of-time compiled Go code.

status
------

2018 Jan 12 update
------
Switch statements now work. They work at the top level
and inside a function.


2018 Jan 11 update
------
With release v0.5.0 the inital test of the `defer`/
`panic`/`recover` mechanism passes. Woot!  There's more
to do here, but the design is solid so filling in
should be quick.

For a stack-unwinding `panic`, we use
what Lua offers, the `error`
mechanic -- to throw -- combined with the
`xpcall` mechanic to catch.

The only limitation I found here is on recursive `xpcall`: if
you are in a panic stack unwind, and then in a defer function,
and your code causes a second error that is *not* a deliberate panic,
then that error will be caught but recover won't return
that error value to the caller of recover. This is due to a wierd
corner case in the implimentation of LuaJIT where
it does not like recursive `xpcall` invocations, and
reports "error in error handling".

I asked on the Lua and LuaJIT mailing lists about
this, and posted on Stack Overflow. So far no
replies. https://stackoverflow.com/questions/48202338/on-latest-luajit-2-1-0-beta3-is-recursive-xpcall-possible

It's a fairly minor limitation, and easy to work
around once you notice the bug: just call `panic`
directly rather than causing the error. Or don't
cause the error at all(!) Or simply use a
different side-band to pass around the value.
Lots of work arounds. One
could also run on the (generally slower) PUC
Lua 5.1.5 as it doesn't have the limitation.


2018 Jan 10 update
------
With release v0.4.1 we have much improved map support.
In `gi`, maps now work as in Go. Multiple-valued queries
properly return the zero-value for the value-type
when the key is missing, and return the 2nd
value correctly. Nil keys and values are
handled properly. `delete` on a map works as
expected, and maintains the `len` property.

2018 Jan 9 update
------
Functions and methods can now be re-defined at the REPL. The
type checker was relaxed to allow this.

We changed over from one luajit C binding to anther. The
new binding is the same one that LuaR uses, so this
enables LuaR exploration.

2018 Jan 8 update
------
Today we landed multiline support. We evalutate
Go expressions as they are entered, and these
can now span multiple lines. This lifts the
prior limitation that meant that functions
and types needed to be defined all on one line.

This was fun to put together. I used the actual gc front end that
parses regular go code. Since gc is written in Go,
why not leverage it! The advantage is that we know we
are building on correct parsing of the whole language.

Of course minor tweaks had to be made to allow statements and
expressions at global scope. Happily, from our experience
adding the same provisions to GopherJS, we knew these
were relatively minor changes. See the updated
demo transcript below in this readme for a multi-line
function definition taking effect.



2018 Jan 7: latest update
------
Today we acheived passing (light) tests for method definition and invocation!

Also a significant discovery for the object system: Steve Donovan's Luar
provides object exchange both ways between Go -> Lua and Lua -> Go.

That should influence our design of our Go source -> Lua source mapping. If we
map in a way that matches what Luar does when it translates from
Go binary -> Lua binary, then our objects will translate cleanly
into binary Go calls made by reflection.

Even more: Luar provides access to the full Go runtime and channels
via reflection. Nice! We don't have to reinvent the wheel, and we
get to use the high-performance multicore Go scheduler.


earlier summary
-----
Early stages, work in progress. Contribute!

Currently incremental type checking is applied
to all code. Slices are bounds-checked at runtime.
Functions, closures and slices, as well as
basic expressions compile and run. For-loops
including for-range loops compile and run.

Much is left to do: maps, structs, switch,
interfaces, imports. If this is exciting to
you, contribute! Existing TODOs/open issues and polite improvement
suggestions can be found here
https://github.com/go-interpreter/gi/issues

However, because we are bulding on the fantastic
front end provided by (Gopherjs)[https://github.com/gopherjs/gopherjs], and the fantastic
backend provided by (LuaJIT)[http://luajit.org/], progress has been
quite rapid.

# the dream

Go, if it only had a decent REPL, could be a great
language for exploratory data analysis.

# the rationale
Go has big advantages over python, R, and Matlab.
It has good type checking, reasonable compiled performance,
and excellent multicore support.

# the aim

We want to provide one excellent integrated REPL for Go.
Exploratory data analysis should not be hampered
by weak type-checking or hard-to-refactor code,
and performance should not suffer just because
you require interaction with your data.


# of course we need a backend to develop against

Considering possible backends for a
reference implementation,
I compared node.js, chez scheme, otto,
gopher-lua, and luajit.

# luajit did what?

Luajit in particular is an amazing
backend to target. In our quick and
dirty 500x500 random matrix multiplication
benchmark, luajit *beat even statically compiled go*
code by a factor of 3x. Go's time was 360 msec.
Luajit's time was 135 msec. Julia uses an optimized
BLAS library for this task and beats both Go
and luajit by multiplying in 6 msec, but
is too immature and too large to be
a viable embedded target.

# installation

Works on Mac OSX and Linux. I have no idea about windows (theoretically it should work on windows, if one can work out the right compiler flags and have the necessary C compiler installed, etc.)

~~~
$ go get -t -u -v github.com/go-interpreter/gi/...
$ cd $GOPATH/src/github.com/go-interpreter/gi && make
$
$ ... wait for gi build to finish, it builds luajit
$     using C, so it takes ~ 20 seconds to install `gi`.
$
$ gi # start me up (will be in $GOPATH/bin/gi now).

====================
gi: a go interpreter
====================
https://github.com/go-interpreter/gi
Copyright (c) 2018, Jason E. Aten, Ph.D.
License: 3-clause BSD. See the LICENSE file at
https://github.com/go-interpreter/gi/blob/master/LICENSE
====================
  [ gi is an interactive Golang environment,
    also known as a REPL or Read-Eval-Print-Loop ]
  [ type ctrl-d to exit ]
  [ type :help for help ]
  [ gi -h for flag help ]
  [ gi -q to start quietly ]
====================
built: '2018-01-08T23:46:07-0600'
last-git-commit-hash: 'db302d2acb37d3c2ba2a0d376b6f233045928730'
nearest-git-tag: 'v0.3.3'
git-branch: 'master'
go-version: 'go_version_go1.9_darwin/amd64'
luajit-version: 'LuaJIT_2.1.0-beta3_--_Copyright_(C)_2005-2017_Mike_Pall._http://luajit.org/'
==================
using this prelude directory: '/Users/jaten/go/src/github.com/go-interpreter/gi/pkg/compiler'
using these files as prelude: array.lua, map.lua, prelude.lua, slice.lua, struct.lua
gi>

gi> a := []string{"howdy", "gophers!"}

gi> a   // ^^ make data using Go's literals. inspect it by typing the variables name.
slice of length 2 is _giSlice{[0]= howdy, [1]= gophers!, }

gi> a[0]  = "you rock" // data can be changed

gi> a
slice of length 2 is _giSlice{[0]= you rock, [1]= gophers!, }

gi> // the Go type checker helps you quickly catch blunders, at compile time.

gi> a[-1] = "compile-time-out-of-bounds-access" 
oops: 'problem detected during Go static type checking: 'where error? err = '1:3: invalid argument: index -1 (constant of type int) must not be negative''' on input 'a[-1] = "compile-time-out-of-bounds-access" 
'

gi> // runtime bounds checks are compiled in too:

gi> a[100] = "runtime-out-of-bounds-access"
error from Lua vm.Pcall(0,0,0): 'run time error'. supplied lua with: '	_gi_SetRangeCheck(a, 100, "runtime-out-of-bounds-access");'
lua stack:
String : 	 ...rc/github.com/go-interpreter/gi/pkg/compiler/prelude.lua:14: index out of range

gi> // We can define functions:

gi> func myFirstGiFunc(a []string) int {
>>>    for i := range a {
>>>      println("our input is a[",i,"] = ", a[i]) 
>>>    };
>>>    return 43
>>> }
func myFirstGiFunc(a []string) int {

	for i := range a {

		println("our input is a[", i, "] = ", a[i])

	}

	return 43

}
gi> myFirstGiFunc(a)
our input is a[	0	] = 	you rock
our input is a[	1	] = 	gophers!

gi> // ^^ and call them. They are tracing-JIT compiled on the LuaJIT vm.

gi> // more compile time type checking, because it rocks:

gi> b := []int{1,1}

gi> myFirstGiFunc(b)
oops: 'problem detected during Go static type checking: 'where error? err = '1:15: cannot use b (variable of type []int) as []string value in argument to myFirstGiFunc''' on input 'myFirstGiFunc(b)
'
gi>
~~~

# editor support

An emacs mode `gigo.el` can be found in the `emacs/` subdirectory
here https://github.com/go-interpreter/gi/tree/master/emacs/gigo.el

M-x `run-gi-golang` to start the interpreter. Pressing ctrl-n will
step through any file that is in `gi-golang` mode. 

Other editors: please contribute!

# how to contribute: getting started

a) Pick an issue from here, https://github.com/go-interpreter/gi/issues, and add a comment that
you are starting work on that feature. Make a branch for your feature, using `git checkout -b yourFeatureName`.

b) Write a test for your feature. Make sure it fails (the test is red), before
moving on to implementation. Tests are quite short. There are many examples are here,
which show the currently implemented features. Add your test at the end of
the compiler/repl_test.go file.

https://github.com/go-interpreter/gi/blob/master/pkg/compiler/repl_test.go

Then simply implement your feature. These are the main files

https://github.com/go-interpreter/gi/blob/master/pkg/compiler/incr.go

https://github.com/go-interpreter/gi/blob/master/pkg/compiler/package.go

https://github.com/go-interpreter/gi/blob/master/pkg/compiler/translate.go

https://github.com/go-interpreter/gi/blob/master/pkg/compiler/statements.go

https://github.com/go-interpreter/gi/blob/master/pkg/compiler/expressions.go

https://github.com/go-interpreter/gi/blob/master/pkg/compiler/luaUtil.go

You wil find it necessary and useful to add print statements to the code. Do this using the `pp()` function, and feel free to leave those prints in during commit. It's a small matter later to take them out, and while you are adding functionality, the debug prints help immensely. You will see them scattered through the code as I've worked. Just leave them there; the verb.Verbose flag and VerboseVerbose can be used to mute them.

The files above derive from GopherJS which compiles Go into Javascript; whereas `gi` translates Go into
Lua. This makes implementation usually very fast, since mostly it is just
above figuring out how to re-write javascript into Lua. You are typically just
checking the syntax of the source-to-source translation. Sometimes some
Lua support functions will be needed. Add them to a new .lua file in `compile/`
directory.

By default, `gi` looks in `./prelude/` relative to its current directory, and this is symlinked to `pkg/compile/` if you are running in `cmd/gi`. Otherwise use the `-prelude` flag to `gi` to tell it where to find its prelude files. All
.lua files found the prelude directory will be sourced during `gi` startup. The default prelude is the `pkg/compile` directory. These files are required for `gi` to work.

c) When you are done, make sure all the tests are green `go test -v` in the compile/ directory.
Run `go fmt` on your code.

d) submit your pull request! (Rebase against master first, please).

# Lua resources

LuaJIT targets Lua 5.1 with some 5.2 extensions.

a) main web site

https://www.lua.org/

b) Programming in Lua by by Roberto Ierusalimschy, the chief architect of Lua.

1st edition. html format (Lua 5.0) https://www.lua.org/pil/contents.html

2nd edition. pdf format (Lua 5.1) https://doc.lagout.org/programmation/Lua/Programming%20in%20Lua%20Second%20Edition.pdf

c) Lua 5.1 Reference Manual, by R. Ierusalimschy, L. H. de Figueiredo, W. Celes
Lua.org, August 2006 

Lua 5.1 https://www.lua.org/manual/5.1/ 

# origin

Author: Jason E. Aten, Ph.D.

License: 3-clause BSD.

Credits: some code here is dervied from the Go standard
libraries, the Go gc compiler,  and from Richard Musiol's excellent Gopherjs project.
This project and those are licensed under the 3-clause BSD license
found in the LICENSE file. The LuaJIT vm and compiler are statically linked
using CGO, and their MIT license can be found in their sub-directories
and online at http://luajit.org/ and https://github.com/LuaJIT/LuaJIT/blob/master/COPYRIGHT
See the subdirectories of vendored and utilized libraries for their
license details.
