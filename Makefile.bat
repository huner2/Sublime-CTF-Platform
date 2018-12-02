@echo off

if [%1] == []       goto :default
if /I %1 == default goto :default
if /I %1 == deps    goto :deps
if /I %1 == clean   goto :clean
if /I %1 == run     goto :run


:default
    @echo on
    go build -o ctfExec.exe
    @echo off
    goto :eof

:deps
    @echo on
    go get -u github.com/gorilla/mux
	go get -u github.com/flosch/pongo2
	go get -u gopkg.in/ini.v1
	go get -u github.com/lib/pq
    go get -u golang.org/x/crypto/blake2b
    @echo off
    goto :eof

:clean
    @echo on
    del /F ctfExec.exe
    @echo off
    goto :eof

:run
    @echo on
    ctfExec.exe
    @echo off
    goto :eof