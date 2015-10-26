# This Makefile is only useful for building binaries for all platforms.
#
# You can build a binary for your platform using the standard "go build"

APPNAME=$(notdir $(PWD))

dist:
	mkdir -p dist/

clean:
	rm -rf dist/

dist/$(APPNAME)_osx.zip: dist
	rm -f dist/$(APPNAME)
	env GOOS=darwin GOARCH=amd64 go build -o dist/$(APPNAME)
	cd dist && zip $(APPNAME)_osx.zip $(APPNAME)
	rm -f dist/$(APPNAME)
all: dist/$(APPNAME)_osx.zip

dist/$(APPNAME)_win32.zip: dist
	rm -f dist/$(APPNAME).exe
	env GOOS=windows GOARCH=386 go build -o dist/$(APPNAME).exe
	cd dist && zip $(APPNAME)_win32.zip $(APPNAME).exe
	rm -f dist/$(APPNAME).exe
all: dist/$(APPNAME)_win32.zip

dist/$(APPNAME)_win64.zip: dist
	rm -f dist/$(APPNAME).exe
	env GOOS=windows GOARCH=amd64 go build -o dist/$(APPNAME).exe
	cd dist && zip $(APPNAME)_win64.zip $(APPNAME).exe
	rm -f dist/$(APPNAME).exe
all: dist/$(APPNAME)_win64.zip

dist/$(APPNAME)_linux.zip: dist
	rm -f dist/$(APPNAME)
	env GOOS=linux GOARCH=386 go build -o dist/$(APPNAME)
	cd dist && zip $(APPNAME)_linux.zip $(APPNAME)
	rm -f dist/$(APPNAME)
all: dist/$(APPNAME)_linux.zip

dist/$(APPNAME)_linux64.zip: dist
	rm -f dist/$(APPNAME)
	env GOOS=linux GOARCH=amd64 go build -o dist/$(APPNAME)
	cd dist && zip $(APPNAME)_linux64.zip $(APPNAME)
	rm -f dist/$(APPNAME)
all: dist/$(APPNAME)_linux64.zip
