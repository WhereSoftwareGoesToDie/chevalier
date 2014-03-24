install: request_sources check
	go install

request_sources: deps request_sources.go
	go build

deps:
	go get

clean:
	rm -f request_sources

check:
	go test

.PHONY: check
