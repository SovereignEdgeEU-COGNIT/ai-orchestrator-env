all: build
.PHONY: all build

BUILD_IMAGE ?= johan/envserver
PUSH_IMAGE ?= johan/envserver:v0.0.1

VERSION := $(shell git rev-parse --short HEAD)
BUILDTIME := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')

GOLDFLAGS += -X 'main.BuildVersion=$(VERSION)'
GOLDFLAGS += -X 'main.BuildTime=$(BUILDTIME)'

build:
	#@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w $(GOLDFLAGS)" -o ./bin/envcli ./cmd/main.go
	@CGO_ENABLED=0 go build -ldflags="-s -w $(GOLDFLAGS)" -o ./bin/envcli ./cmd/main.go

container:
	docker build -t $(BUILD_IMAGE) .

push:
	docker tag $(BUILD_IMAGE) $(PUSH_IMAGE) 
	docker push $(BUILD_IMAGE)
	docker push $(PUSH_IMAGE)

test:
	@cd pkg/core; go test -v --race
	@cd pkg/database; go test -v --race
	@cd pkg/server; go test -v --race

install:
	cp ./bin/envcli /usr/local/bin

startdb: 
	docker run -d -p 5432:5432 -e POSTGRES_PASSWORD=rFcLGNkgsNtksg6Pgtn9CumL4xXBQ7 --restart unless-stopped timescale/timescaledb:latest-pg16
