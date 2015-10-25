# gow
Watching directory & Running command written Go.  
Easy, Portability and Fastest.

## Naming:
**go** **w**atch and run.  
Why not inculde **r**un?  
`gow` is easy typing. `gowr` is hard to type in dvorak.

## Install:
```bash
go get -u github.com/zchee/gow
```

## Usage:
```bash
Usage of ./gow:
  -path string
        Watch directory path (default "./")
```

## Requirements:
- Go
- [WIP] OS X only
  - Now, `gow` use [`Apple OS X FSEvents`](https://developer.apple.com/library/mac/documentation/Darwin/Reference/FSEvents_Ref/) and [`go-fsnotify/fsevents`](https://github.com/go-fsnotify/fsevents)
  - Support `linux` and `Windows` use [go-fsnotify/fsnotify](https://github.com/go-fsnotify/fsnotify) if I feel like it :)

## Goals:
- Support multi platform use Go binary build.
- An alternative all automation build tools. Bye `guard`, `grunt`, `gulp`.
- Pluggable parse file interface.

## Support parse files:

| File        | State | Language     |
|-------------|-------|--------------|
| flag        | Yes   | Shell        |
| `gulp.js`   | No    | JavaScript   |
| `grunt.js`  | No    | JavaScript   |
| `Guardfile` | No    | Guard DSL    |
| `Makefile`  | No    | Makefile DSL |
| etc...      | No    | Any          |

## License:
The MIT License (MIT)

## Author:
Copyright (c) 2015- zchee aka Koichi Shiraishi
