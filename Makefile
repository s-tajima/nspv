deps:
	go get -v -t -d ./...
	go get -u github.com/schrej/godacov


build:
	go build -v .

test:
	go test -v -coverprofile=coverage.out .

godacov:
	godacov -t $${CODACY_TOKEN} -r ./coverage.out -c $$(git rev-parse HEAD)
