build:
	@go build -o bin/eardis

run: build
	@./bin/eardis

test:
	@go test -v ./...
