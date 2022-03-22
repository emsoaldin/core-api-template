
GO ?= go
GOFMT ?= gofmt "-s"
PACKAGES ?= $(shell $(GO) list ./... | grep -v /vendor/)
GOFILES := find . -name "*.go" -type f -not -path "./vendor/*"
COVERFILE := coverage.txt
SERVICE_NAME := core-api
VERSION := `date -u +1.%Y%m%d.%H%M%S`

.PHONY: deps
deps:
	$(GO) mod download

.PHONY: clean
clean:
	$(GO) clean -modcache -cache -i
	rm -rf ./build

.PHONY: fmt
fmt:
	$(GOFILES) | xargs $(GOFMT) -w

.PHONY: fmt-check
fmt-check:
	@files=$$($(GOFILES) | xargs $(GOFMT) -l); if [ -n "$$files" ]; then \
		echo "Please run 'make fmt' and commit the result:"; \
		echo "$${files}"; \
		exit 1; \
		fi;

.PHONY: test
test: fmt-check
	$(GO) test -v -coverprofile=$(COVERFILE) ./...

coverage: test
	$(GO) tool cover -html=$(COVERFILE)


.PHONY: swagger
swagger:
	rm -rf ./docs
	swag init --parseDependency --parseInternal --parseDepth 1 -generalInfo cmd/$(SERVICE_NAME)/main.go

.PHONY: build
build: 
	CGO_ENABLED=0 $(GO) build -installsuffix 'static' -ldflags "-X api/version.API=${VERSION}" -o ./build/$(SERVICE_NAME) ./cmd/$(SERVICE_NAME)

.PHONY: run
run:
	make build
	env $$(cat .env) ./build/$(SERVICE_NAME)
