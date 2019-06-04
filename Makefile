# This Makefile is meant to be used by people that do not usually work
# with Go source code. If you know what GOPATH is then you probably
# don't need to bother with make.

.PHONY: graviton android ios graviton-cross swarm svm all test clean
.PHONY: graviton-linux graviton-linux-386 graviton-linux-amd64 graviton-linux-mips64 graviton-linux-mips64le
.PHONY: graviton-linux-arm graviton-linux-arm-5 graviton-linux-arm-6 graviton-linux-arm-7 graviton-linux-arm64
.PHONY: graviton-darwin graviton-darwin-386 graviton-darwin-amd64
.PHONY: graviton-windows graviton-windows-386 graviton-windows-amd64

GOBIN = $(shell pwd)/build/bin
GO ?= latest

graviton:
	build/env.sh go run build/ci.go install ./cmd/graviton
	@echo "Done building."
	@echo "Run \"$(GOBIN)/graviton\" to launch graviton."

swarm:
	build/env.sh go run build/ci.go install ./cmd/swarm
	@echo "Done building."
	@echo "Run \"$(GOBIN)/swarm\" to launch swarm."

all:
	build/env.sh go run build/ci.go install

android:
	build/env.sh go run build/ci.go aar --local
	@echo "Done building."
	@echo "Import \"$(GOBIN)/graviton.aar\" to use the library."

ios:
	build/env.sh go run build/ci.go xcode --local
	@echo "Done building."
	@echo "Import \"$(GOBIN)/Graviton.framework\" to use the library."

test: all
	build/env.sh go run build/ci.go test

lint: ## Run linters.
	build/env.sh go run build/ci.go lint

clean:
	./build/clean_go_build_cache.sh
	rm -fr build/_workspace/pkg/ $(GOBIN)/*

# The devtools target installs tools required for 'go generate'.
# You need to put $GOBIN (or $GOPATH/bin) in your PATH to use 'go generate'.

devtools:
	env GOBIN= go get -u golang.org/x/tools/cmd/stringer
	env GOBIN= go get -u github.com/kevinburke/go-bindata/go-bindata
	env GOBIN= go get -u github.com/fjl/gencodec
	env GOBIN= go get -u github.com/golang/protobuf/protoc-gen-go
	env GOBIN= go install ./cmd/abigen
	@type "npm" 2> /dev/null || echo 'Please install node.js and npm'
	@type "solc" 2> /dev/null || echo 'Please install solc'
	@type "protoc" 2> /dev/null || echo 'Please install protoc'

swarm-devtools:
	env GOBIN= go install ./cmd/swarm/mimegen

# Cross Compilation Targets (xgo)

graviton-cross: graviton-linux graviton-darwin graviton-windows graviton-android graviton-ios
	@echo "Full cross compilation done:"
	@ls -ld $(GOBIN)/graviton-*

graviton-linux: graviton-linux-386 graviton-linux-amd64 graviton-linux-arm graviton-linux-mips64 graviton-linux-mips64le
	@echo "Linux cross compilation done:"
	@ls -ld $(GOBIN)/graviton-linux-*

graviton-linux-386:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/386 -v ./cmd/graviton
	@echo "Linux 386 cross compilation done:"
	@ls -ld $(GOBIN)/graviton-linux-* | grep 386

graviton-linux-amd64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/amd64 -v ./cmd/graviton
	@echo "Linux amd64 cross compilation done:"
	@ls -ld $(GOBIN)/graviton-linux-* | grep amd64

graviton-linux-arm: graviton-linux-arm-5 graviton-linux-arm-6 graviton-linux-arm-7 graviton-linux-arm64
	@echo "Linux ARM cross compilation done:"
	@ls -ld $(GOBIN)/graviton-linux-* | grep arm

graviton-linux-arm-5:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm-5 -v ./cmd/graviton
	@echo "Linux ARMv5 cross compilation done:"
	@ls -ld $(GOBIN)/graviton-linux-* | grep arm-5

graviton-linux-arm-6:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm-6 -v ./cmd/graviton
	@echo "Linux ARMv6 cross compilation done:"
	@ls -ld $(GOBIN)/graviton-linux-* | grep arm-6

graviton-linux-arm-7:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm-7 -v ./cmd/graviton
	@echo "Linux ARMv7 cross compilation done:"
	@ls -ld $(GOBIN)/graviton-linux-* | grep arm-7

graviton-linux-arm64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm64 -v ./cmd/graviton
	@echo "Linux ARM64 cross compilation done:"
	@ls -ld $(GOBIN)/graviton-linux-* | grep arm64

graviton-linux-mips:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips --ldflags '-extldflags "-static"' -v ./cmd/graviton
	@echo "Linux MIPS cross compilation done:"
	@ls -ld $(GOBIN)/graviton-linux-* | grep mips

graviton-linux-mipsle:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mipsle --ldflags '-extldflags "-static"' -v ./cmd/graviton
	@echo "Linux MIPSle cross compilation done:"
	@ls -ld $(GOBIN)/graviton-linux-* | grep mipsle

graviton-linux-mips64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips64 --ldflags '-extldflags "-static"' -v ./cmd/graviton
	@echo "Linux MIPS64 cross compilation done:"
	@ls -ld $(GOBIN)/graviton-linux-* | grep mips64

graviton-linux-mips64le:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips64le --ldflags '-extldflags "-static"' -v ./cmd/graviton
	@echo "Linux MIPS64le cross compilation done:"
	@ls -ld $(GOBIN)/graviton-linux-* | grep mips64le

graviton-darwin: graviton-darwin-386 graviton-darwin-amd64
	@echo "Darwin cross compilation done:"
	@ls -ld $(GOBIN)/graviton-darwin-*

graviton-darwin-386:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=darwin/386 -v ./cmd/graviton
	@echo "Darwin 386 cross compilation done:"
	@ls -ld $(GOBIN)/graviton-darwin-* | grep 386

graviton-darwin-amd64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=darwin/amd64 -v ./cmd/graviton
	@echo "Darwin amd64 cross compilation done:"
	@ls -ld $(GOBIN)/graviton-darwin-* | grep amd64

graviton-windows: graviton-windows-386 graviton-windows-amd64
	@echo "Windows cross compilation done:"
	@ls -ld $(GOBIN)/graviton-windows-*

graviton-windows-386:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=windows/386 -v ./cmd/graviton
	@echo "Windows 386 cross compilation done:"
	@ls -ld $(GOBIN)/graviton-windows-* | grep 386

graviton-windows-amd64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=windows/amd64 -v ./cmd/graviton
	@echo "Windows amd64 cross compilation done:"
	@ls -ld $(GOBIN)/graviton-windows-* | grep amd64
