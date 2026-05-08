.PHONY: build test 

BUILD_CMD = go build -o bin/reg.exe cmd/reg/main.go
TEST_CMD = go test ./... -v

build:
	$(BUILD_CMD)

test:
	$(TEST_CMD)
