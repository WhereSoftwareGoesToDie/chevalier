install: index_sources check
	go install

index_sources: deps index_sources.go
	go build

deps:
	go get

clean:
	rm -f index_sources

check:
	go test

.PHONY: check
