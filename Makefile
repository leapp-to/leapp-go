.PHONY: test
test:
	@go clean
	@go test ./...

.PHONY: all build
all build:
	@go build -o build/actor-stdout cmd/actor-stdout/main.go 
	@go build -o build/leapp-daemon cmd/leapp-daemon/main.go 

.PHONY: clean
clean:
	@rm -rf build/
