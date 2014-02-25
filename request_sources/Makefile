request_sources: install

install: build check
	go install

build: deps
	go build

deps:
	go get

clean:
	rm -f request_sources

check:
	go test
