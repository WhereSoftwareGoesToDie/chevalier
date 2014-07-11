all: install chevalierd request_sources check_chevalier

install: build check
	go get
	go install

chevalierd: 
	cd chevalierd ; make

request_sources: 
	cd request_sources; make

check_chevalier: 
	cd check_chevalier; make

build: protobuf 
	go build

protobuf: goprotobuf
	cd protobuf ;  protoc --go_out=.. *.proto

goprotobuf:
	go get code.google.com/p/goprotobuf/proto
	go get code.google.com/p/goprotobuf/protoc-gen-go

check: protobuf
	go test

.PHONY : all
.PHONY : protobuf
.PHONY : request_sources
.PHONY : chevalierd
