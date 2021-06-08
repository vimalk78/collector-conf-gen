

.PHONY: force

EXE=collector-conf-gen

build: bin/$(EXE)

bin/$(EXE): force
	go build $(BUILD_OPTS) -o $@ ./cmd/main

generate: build
	./bin/collector-conf-gen
