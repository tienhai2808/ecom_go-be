ROOT := .
TMP_DIR := tmp
BIN := $(TMP_DIR)/main.exe
TESTDATA_DIR := testdata

.PHONY: build run clean

build:
	@echo "Building..."
	@if not exist $(TMP_DIR) mkdir $(TMP_DIR)
	go build -o $(BIN) ./cmd/api

run: build
	@echo "Running..."
	@$(BIN)

clean:
	@echo "Cleaning..."
	@rm -rf $(TMP_DIR)