.PHONY: build cluster test

default: build

build:
	go build -o bin/swarm    github.com/susy-go/susy-graviton/cmd/swarm
	go build -o bin/graviton     github.com/susy-go/susy-graviton/cmd/graviton
	go build -o bin/bootnode github.com/susy-go/susy-graviton/cmd/bootnode

cluster: build
	scripts/boot-cluster.sh

test:
	go test -v github.com/susy-go/susy-graviton/swarm/...
