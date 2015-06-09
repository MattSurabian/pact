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