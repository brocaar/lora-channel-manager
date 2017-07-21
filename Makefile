.PHONY: build clean test package package-deb
PKGS := $(shell go list ./... | grep -v /vendor/ | grep -v loraserver/api | grep -v /migrations | grep -v /static)
VERSION := $(shell git describe --always)
GOOS ?= linux
GOARCH ?= amd64

build:
	@echo "Compiling source for $(GOOS) $(GOARCH)"
	@mkdir -p build
	@GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags "-X main.version=$(VERSION)" -o build/lora-gateway-config$(BINEXT) cmd/lora-gateway-config/main.go

clean:
	@echo "Cleaning up workspace"
	@rm -rf build
	@rm -rf docs/public

test:
	@echo "Running tests"
	@for pkg in $(PKGS) ; do \
		golint $$pkg ; \
	done
	@go vet $(PKGS)
	@go test -p 1 -v $(PKGS)

package: clean build
	@echo "Creating package for $(GOOS) $(GOARCH)"
	@mkdir -p dist/tar/$(VERSION)
	@cp build/* dist/tar/$(VERSION)
	@cd dist/tar/$(VERSION) && tar -pczf ../lora_gateway_config_$(VERSION)_$(GOOS)_$(GOARCH).tar.gz .
	@rm -rf dist/tar/$(VERSION)

package-deb:
	@cd packaging && TARGET=deb ./package.sh
