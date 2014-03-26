install: chevalier_status check
	go install

chevalier_status: deps chevalier_status.go
	go build

deps:
	go get

clean:
	rm -f chevalier_status

check:
	go test

.PHONY: check
