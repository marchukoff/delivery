format:
	go fmt ./...
.PHONY: format

test: format
	go test -v -count=1 ./...
.PHONY: test