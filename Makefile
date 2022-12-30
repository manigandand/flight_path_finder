# Overcome m2c: Command not found
# https://lists.mcs.anl.gov/pipermail/petsc-users/2010-November/007299.html
 %.o: %.mod
.PHONY: build

# Makefile
export GO111MODULE=on

GO_CMD=go
GO_MOD_TIDY=$(GO_CMD) mod tidy
GO_BUILD=$(GO_CMD) build
GO_BUILD_RACE=$(GO_CMD) build -race
GO_TEST=$(GO_CMD) test
GO_TEST_VERBOSE=$(GO_CMD) test -v
GO_TEST_COVER=$(GO_CMD) test -cover -count=1 -p=1
GO_INSTALL=$(GO_CMD) install -v

SERVER_BIN=volumefi
SERVER_BIN_DIR=.
SERVER_DIR=.
SERVER_MAIN=main.go

SOURCE_PKG_DIR= .
SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

lint:
	./golangci_linter.sh ./... $(LINT_ENV)

gotests: lint
	$(GO_TEST_VERBOSE) ./... --failfast

test: gotests

build:
	$(GO_MOD_TIDY)
	env GOOS=linux GARCH=amd64 CGO_ENABLED=0 $(GO_BUILD) -o volumefi -ldflags="-X 'main.version=$(version)'" .
	@chmod 755 $(SERVER_BIN)

run_clean: test build
	./$(SERVER_BIN)

run: build
	./$(SERVER_BIN)