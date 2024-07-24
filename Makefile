.SILENT:
.DEFAULT_GOAL := run

BINARY_NAME=aitime

run: build
	./bin/$(BINARY_NAME)-linux

build:
	GOARCH=amd64 GOOS=linux go build -o bin/$(BINARY_NAME)-linux .
	GOARCH=amd64 GOOS=darwin go build -o bin/$(BINARY_NAME)-darwin .
	GOARCH=amd64 GOOS=windows go build -o bin/$(BINARY_NAME)-windows .

clean:
	rm -rf bin/
	rm -rf cover.out
	rm -rf cover.html

test:
	go test ./...

testcoverage:
	go test -coverprofile cover.out ./...
	go tool cover -html cover.out -o cover.html
