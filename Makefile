# Variables
GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
BUILD_DIR := ./build
GIT_COMMIT := $(shell git rev-list -1 HEAD)
DEBUG_FLAGS := -ldflags "-X main.GitCommit=$(GIT_COMMIT)" -gcflags "all=-N -l"
RELEASE_FLAGS := -ldflags="-w -s -X main.GitCommit=$(GIT_COMMIT)" -gcflags="all=-trimpath=${PWD}"
BUILD_CMD := $(GOBUILD) -o $(BUILD_DIR)/ -v
OS := linux
SkipTest := false

# Default target
all: clean test build

# Clean the build directory
clean:
	@echo "[+] Clean the build directory"
	@$(GOCLEAN)
	@rm -rf $(BUILD_DIR)

# Build the project
build: clean test
	@echo "[+] Building project"
	@$(BUILD_CMD) $(DEBUG_FLAGS)

# build a release binary
release: clean test
	@echo "[+] Building a release binary"
	@GOOS=$(OS) $(BUILD_CMD) $(RELEASE_FLAGS)

# build a debug binary
debug: clean test
	@echo "[+] Building a debug binary"
	@GOOS=$(OS) $(BUILD_CMD) $(DEBUG_FLAGS)

# Run tests
test:
ifeq ($(SkipTest), false)
	@echo "[+] Running tests"
	@$(GOTEST) -v ./... -count=1 -coverprofile=coverage.out
else
	@echo "[+] Skipping tests"
endif

cover: test
	@go tool cover -html=coverage.out

# Specify the phony targets
.PHONY: all clean build test cover
