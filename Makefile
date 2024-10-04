.PHONY: lint test vendor clean

export GO123MODULE=on

default: lint test

lint:
	golangci-lint version
	golangci-lint run

test:
	go test -v -cover ./...

yaegi_test:
	yaegi test -v .

vendor:
	go mod vendor

clean:
	rm -rf ./vendor