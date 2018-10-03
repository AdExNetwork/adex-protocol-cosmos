all: build

VERSION="1.0.0"
COMMIT=$(shell git rev-parse HEAD)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)

LDFLAGS = -ldflags "-X main.Version=${VERSION} -X main.Commit=${COMMIT} -X main.Branch=${BRANCH}"

cli:
	go build ${LDFLAGS} -o ./adexcli cmd/adexcli/main.go

daemon:
	go build ${LDFLAGS} -o ./adexd cmd/adexd/main.go

build: cli daemon
#test:
#	go test -v $(PACKAGES)

#silent-test:
#	go test $(PACKAGES)

deps:
	go get

format:
	go fmt ./...
