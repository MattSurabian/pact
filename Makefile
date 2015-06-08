GO ?= go
GOPATH := $(CURDIR)/_vendor:$(GOPATH)

default: build

build:
	./vendor.sh
	cd $(CURDIR) && $(GO) test ./... && $(GO) build