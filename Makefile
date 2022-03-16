all: build-run

run-main:
	go run cmd/fsqm/*.go

clear:
	rm -rf bin

build:
	go build -o bin/fsqm cmd/fsqm/*.go

install:
	go install cmd/fsqm/*.go

snapshot: clear
	sh ./scripts/docker-deploy.sh

build-run: clear build
	./bin/fsqm
