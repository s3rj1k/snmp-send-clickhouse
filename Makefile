GO_BIN ?= go
CURL_BIN ?= curl
SHELL_BIN ?= sh
BIN_NAME = snmpmon

all: clean build

build:
	$(GO_BIN) build -ldflags="-s -w" -o $(BIN_NAME) -v

clean:
	$(GO_BIN) clean
	rm -f $(BIN_NAME)

deps: check-gopath
	$(GO_BIN) get -u gopkg.in/go-playground/validator.v9
	$(GO_BIN) get -u github.com/mailru/go-clickhouse

	$(CURL_BIN) -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | $(SHELL_BIN) -s -- -b ${GOPATH}/bin v1.15.0
	$(GO_BIN) get -u github.com/Quasilyte/go-consistent

test:
	$(GO_BIN) test -failfast ./...

lint:
	golangci-lint run
	go-consistent -v ./...

check-gopath:
ifndef GOPATH
	$(error GOPATH is undefined)
endif