# This Makefile is only useful for building binaries for all platforms.
#
# You can build a binary for your platform using the standard "go build"

APPNAME=$(basename $(PWD))

dist:
	mkdir -p dist/

clean:
	rm -rf dist/

dist/$(APPNAME)_osx: dist
	env GOOS=darwin GOARCH=amd64 go build -o dist/$(APPNAME)_osx
all: dist/$(APPNAME)_osx

dist/$(APPNAME)_win32: dist
	env GOOS=windows GOARCH=386 go build -o dist/$(APPNAME)_win32.exe
all: dist/$(APPNAME)_win32

dist/$(APPNAME)_win64: dist
	env GOOS=windows GOARCH=amd64 go build -o dist/$(APPNAME)_win64.exe
all: dist/$(APPNAME)_win64

dist/$(APPNAME)_linux: dist
	env GOOS=linux GOARCH=386 go build -o dist/$(APPNAME)_linux
all: dist/$(APPNAME)_linux

dist/$(APPNAME)_linux64: dist
	env GOOS=linux GOARCH=amd64 go build -o dist/$(APPNAME)_linux64
all: dist/$(APPNAME)_linux64
