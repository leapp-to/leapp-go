 BAD_FMT_FILES:= $(shell find . -iname '*.go' | grep -v /vendor/ | xargs gofmt -s -l)
 BAD_IMP_FILES:= $(shell goimports -l .)

 ROOT_PATH=$(PREFIX)/usr/libexec/leapp

.PHONY: fmt
fmt:
	@test -z $(BAD_FMT_FILES) || (echo -e "gofmt failed in the following file(s):\n$(BAD_FMT_FILES)" && exit 1)

.PHONY: imports
imports:
	@test -z $(BAD_IMP_FILES) || (echo -e "goimports detected problems in the following file(s):\n$(BAD_IMP_FILES)" && exit 1)

.PHONY: lint
lint:
	@golint -set_exit_status ./... 

.PHONY: vet
vet:
	@go vet ./...

.PHONY: test
test:
	@go clean
	@sqlite3 /tmp/actors.db < res/audit-layout.sql
	@LEAPP_STORE_PATH=/tmp/actors.db LEAPP_ACTOR_API=/tmp/actor-api.sock go test ./...

.PHONY: test-all
test-all: fmt lint vet imports test

.PHONY: all build
all build:
	@go build -o build/actor-stdout cmd/actor-stdout/main.go 
	@go build -o build/leapp-daemon cmd/leapp-daemon/main.go 

.PHONY: install-deps
install-deps:
	@go get -t -v ./...

.PHONY: install
install:
	install -Dd $(ROOT_PATH) 
	install -m0755 build/leapp-daemon $(ROOT_PATH) 
	install -m0755 build/actor-stdout $(ROOT_PATH)

.PHONY: clean
clean:
	@rm -rf build/
