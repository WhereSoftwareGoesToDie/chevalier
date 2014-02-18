all: install chevalier

install: build check
	go install

chevalier:
	cd chevalier ; make

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
