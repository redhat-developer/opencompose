all: build



build:
    go build main.go

install:
    go install main.go

test:
	go test ./...
