strip_frame_sources: install

install: build check
	go install

build: deps
	go build

deps:
	go get

clean:
	rm -f strip_frame_sources

check:
	go test
