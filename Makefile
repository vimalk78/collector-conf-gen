.PHONY: force

EXE=collector-conf-gen

build: bin/$(EXE)

bin/$(EXE): force
	go build $(BUILD_OPTS) -o $@ ./cmd/main

generate: build
	./bin/collector-conf-gen

test:
	go test -cover ./internal/... $(COVER)

cover: COVER=-coverprofile=coverage.out
cover: test
	go tool cover -html=coverage.out
