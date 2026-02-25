# Proto version to use from mozilla-ai/mcpd-proto repository.
PROTO_VERSION := v0.1.0

# Base URL for downloading proto files from GitHub.
PROTO_BASE_URL := https://raw.githubusercontent.com/mozilla-ai/mcpd-proto/$(PROTO_VERSION)

# Directories.
TMP_DIR := tmp
PROTO_TMP_DIR := $(TMP_DIR)/plugins
OUT_DIR := pkg/plugins/v1

# Proto files to download.
PROTO_FILES := plugin.proto

.PHONY: all
all: fetch-protos generate

.PHONY: fetch-protos
fetch-protos:
	@echo "Fetching proto files from mcpd-proto $(PROTO_VERSION)..."
	@mkdir -p $(PROTO_TMP_DIR)
	@for file in $(PROTO_FILES); do \
		echo "  Downloading $$file..."; \
		curl -sSfL $(PROTO_BASE_URL)/plugins/v1/$$file -o $(PROTO_TMP_DIR)/$$file; \
	done
	@echo "Proto files downloaded successfully."

.PHONY: generate
generate:
	@echo "Generating Go code from proto files..."
	@mkdir -p $(OUT_DIR)
	@protoc \
		--proto_path=$(PROTO_TMP_DIR) \
		--go_out=$(OUT_DIR) \
		--go_opt=paths=source_relative \
		--go-grpc_out=$(OUT_DIR) \
		--go-grpc_opt=paths=source_relative \
		$(PROTO_TMP_DIR)/*.proto
	@echo "Code generation complete."

.PHONY: clean
clean:
	@echo "Cleaning generated files and temporary directories..."
	@rm -rf $(TMP_DIR)
	@find $(OUT_DIR) -name "*.pb.go" -delete 2>/dev/null || true
	@echo "Clean complete."

.PHONY: update-proto-version
update-proto-version:
	@echo "Current proto version: $(PROTO_VERSION)"
	@echo "To update, modify PROTO_VERSION in this Makefile, then run 'make clean all'"

.PHONY: lint
lint:
	@echo "Running linter..."
	golangci-lint run --fix -v

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all                  - Fetch protos and generate code (default)"
	@echo "  fetch-protos         - Download proto files from mcpd-proto"
	@echo "  generate             - Generate Go code from downloaded protos"
	@echo "  clean                - Remove generated files and temp directories"
	@echo "  lint                 - Run golangci-lint with auto-fix"
	@echo "  update-proto-version - Show current proto version"
	@echo "  help                 - Show this help message"
