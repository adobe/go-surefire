local-build:
	go clean -i ./...
	go get -v ./...
	go test -v ./...
	go build --v ./...

.PHONY: local-build
