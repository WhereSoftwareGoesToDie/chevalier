all: install

install: build check
	go install

build: deps
	go build

check: 
	go test

deps:
	go get

.PHONY : all
