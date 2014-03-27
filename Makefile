all: install chevalierd index_sources request_sources strip_frame_sources check_chevalier

install: build check
	go install

chevalierd: 
	cd chevalierd ; make

index_sources: 
	cd index_sources; make

request_sources: 
	cd request_sources; make

strip_frame_sources: 
	cd strip_frame_sources; make

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
.PHONY : index_sources
.PHONY : strip_frame_sources
