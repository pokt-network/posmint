COMMIT := $(shell git log -1 --format='%H')
PACKAGES=$(shell go list ./...)
LEDGER_ENABLED ?= true
BINDIR ?= $(GOPATH)/bin

export GO111MODULE = on

all: lint test build

########################################
### Build

build: go.sum
	@go build -mod=readonly ./...
.PHONY: build

update-swagger-docs: statik
	$(BINDIR)/statik -src=client/lcd/swagger-ui -dest=client/lcd -f -m
	@if [ -n "$(git status --porcelain)" ]; then \
        echo "\033[91mSwagger docs are out of sync!!!\033[0m";\
        exit 1;\
    else \
    	echo "\033[92mSwagger docs are in sync\033[0m";\
    fi
.PHONY: update-swagger-docs

########################################
### Tools & dependencies

go-mod-cache: go.sum
	@echo "--> Download go modules to local cache"
	@go mod download
.PHONY: go-mod-cache

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	@go mod verify
	@go mod tidy

distclean:
	rm -rf \
    gitian-build-darwin/ \
    gitian-build-linux/ \
    gitian-build-windows/ \
    .gitian-builder-cache/
.PHONY: distclean

########################################
### Documentation

godocs:
	@echo "--> Wait a few seconds and visit http://localhost:6060/pkg/github.com/pokt-network/posmint/types"
	godoc -http=:6060

########################################
### Testing

test: test_unit

test_ledger_mock:
	@go test -mod=readonly `go list github.com/pokt-network/posmint/crypto` -tags='cgo ledger test_ledger_mock'

test_ledger: test_ledger_mock
	@go test -mod=readonly -v `go list github.com/pokt-network/posmint/crypto` -tags='cgo ledger'

test_unit:
	@go test -mod=readonly $(PACKAGES) -tags='ledger test_ledger_mock'

test_race:
	@go test -mod=readonly -race $(PACKAGES)

.PHONY: test 

test_cover:
	@bash -x tests/test_cover.sh

lint:	
	@go vet ./...
	@go fmt ./...
	@go mod verify
.PHONY: lint

benchmark:
	@go test -mod=readonly -bench=. $(PACKAGES)
.PHONY: benchmark

########################################
### Devdoc

DEVDOC_SAVE = docker commit `docker ps -a -n 1 -q` devdoc:local

devdoc_init:
	docker run -it -v "$(CURDIR):/go/src/github.com/pokt-network/posmint" -w "/go/src/github.com/pokt-network/posmint" tendermint/devdoc echo
	$(call DEVDOC_SAVE)

devdoc:
	docker run -it -v "$(CURDIR):/go/src/github.com/pokt-network/posmint" -w "/go/src/github.com/pokt-network/posmint" devdoc:local bash

devdoc_clean:
	docker rmi -f $$(docker images -f "dangling=true" -q)

devdoc_update:
	docker pull tendermint/devdoc

.PHONY: devdoc devdoc_clean devdoc_init devdoc_save devdoc_update
