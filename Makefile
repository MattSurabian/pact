GO_CMD=$(shell which go 2>/dev/null)
GO_TEST=$(GO_CMD) test
GO_BUILD=$(GO_CMD) build
GOPATH := $(CURDIR)/_vendor:$(GOPATH)

default: prereq build

prereq:
ifeq ($(GO_CMD),)
	$(error "ERROR: Go is not installed, check out the Golang site: https://golang.org/doc/install or your OS's package manager for information on how to install it.")
endif

build:
	./vendor.sh
	cd $(CURDIR) && $(GO_TEST) ./... && $(GO_BUILD)

release: pre_release build_amd64_win build_386_win build_amd64_mac build_386_mac build_amd64_linux build_386_linux

pre_release:
	./vendor.sh
	cd $(CURDIR) && $(GO_TEST) ./...
	test -d release || mkdir release

build_amd64_win:
	env GOOS=windows GOARCH=amd64 $(GO_BUILD) -o release/pact.exe
	cd release && zip windows-amd64-pact.zip pact.exe && rm pact.exe

build_386_win:
	env GOOS=windows GOARCH=386 $(GO_BUILD) -o release/pact.exe
	cd release && zip windows-i386-pact.zip pact.exe && rm pact.exe

build_amd64_mac:
	env GOOS=darwin GOARCH=amd64 $(GO_BUILD) -o release/pact
	cd release && zip mac-amd64-pact.zip pact && rm pact

build_386_mac:
	env GOOS=darwin GOARCH=386 $(GO_BUILD) -o release/pact
	cd release && zip mac-i386-pact.zip pact && rm pact

build_amd64_linux:
	env GOOS=linux GOARCH=amd64 $(GO_BUILD) -o release/pact
	cd release && zip linux-amd64-pact.zip pact && rm pact

build_386_linux:
	env GOOS=linux GOARCH=386 $(GO_BUILD) -o release/pact
	cd release && zip linux-i386-pact.zip pact && rm pact