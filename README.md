gow
===
Watching directory & Running command written Go.  
Easy, Portability(?) and Fastest.

Install:
--------
```bash
go get -u github.com/zchee/gow
```

Usage:
------
```bash
Usage of ./gow:
  -path string
        Watch directory path (default "./")
```

Requirements:
-------------
- Go
- [WIP] OS X only
  - Now, `gow` use [`Apple OS X FSEvents`](https://developer.apple.com/library/mac/documentation/Darwin/Reference/FSEvents_Ref/) and [`go-fsnotify/fsevents`](https://github.com/go-fsnotify/fsevents)
  - Support `linux` and `Windows` use [go-fsnotify/fsnotify](https://github.com/go-fsnotify/fsnotify) if I feel like it :)

Goals:
------
- Support multi platform use Go binary build.
- An alternative all automation build tools. Bye `guard`, `grunt`, `gulp`.
- Pluggable parse file interface.
- Parse...
  - [ ] `gulp.js`
  - [ ] `grunt.js`
  - [ ] `Guardfile`
  - [ ] `Makefile`
  - [ ] etc...

License:
--------
The MIT License (MIT)

Author:
-------
Copyright (c) 2015 zchee aka Koichi Shiraishi
