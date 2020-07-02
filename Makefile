deps:
	go get -v -t -d ./...

build:
	go build -v .

test:
	go test -v .
