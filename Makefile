all: build-run

run-main:
	go run cmd/fsqm/*.go

clear:
	rm -rf bin

build:
	go build -o bin/fsqm cmd/fsqm/*.go

clear-build: clear build

install: clear-build
	go install cmd/fsqm/*.go

snapshot:
	cd scripts && sudo chmod +x ./docker-deploy.sh
	cd scripts && sudo bash ./docker-deploy.sh

build-run: clear build
	./bin/fsqm
