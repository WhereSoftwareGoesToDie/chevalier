install: check_chevalier check
	go install

check_chevalier: deps check_chevalier.go
	go build

deps:
	go get

clean:
	rm -f check_chevalier

check:
	go test

.PHONY: check
