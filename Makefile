# VictoriaLogs MCP Server
# ========================

# 變數
BINARY_NAME := vlmcp
BINARY_DIR := bin
CMD_DIR := cmd/vlmcp
GO := go
GOFLAGS := -v
LDFLAGS := -ldflags "-s -w -X github.com/vincent119/victorialogs-mcp/pkg/version.Version=$(shell git describe --tags --always --dirty 2>/dev/null || echo 'dev')"

# 預設目標
.DEFAULT_GOAL := help

# ============================================================
# 建置
# ============================================================

.PHONY: build
build: ## 建置 binary
	@mkdir -p $(BINARY_DIR)
	$(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BINARY_DIR)/$(BINARY_NAME) ./$(CMD_DIR)

.PHONY: build-linux
build-linux: ## 建置 Linux binary
	@mkdir -p $(BINARY_DIR)
	GOOS=linux GOARCH=amd64 $(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BINARY_DIR)/$(BINARY_NAME)-linux-amd64 ./$(CMD_DIR)

.PHONY: build-all
build-all: build build-linux ## 建置所有平台

# ============================================================
# 開發
# ============================================================

.PHONY: run
run: ## 執行開發版本
	$(GO) run ./$(CMD_DIR) --config configs/config.yaml

.PHONY: dev
dev: ## 開發模式（使用 air 熱重載，需先安裝 air）
	@which air > /dev/null || (echo "請先安裝 air: go install github.com/cosmtrek/air@latest" && exit 1)
	air

# ============================================================
# 測試
# ============================================================

.PHONY: test
test: ## 執行測試
	$(GO) test -race -cover ./...

.PHONY: test-verbose
test-verbose: ## 執行測試（詳細輸出）
	$(GO) test -race -cover -v ./...

.PHONY: test-coverage
test-coverage: ## 產生測試覆蓋率報告
	$(GO) test -race -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "覆蓋率報告已產生: coverage.html"

# ============================================================
# 程式碼品質
# ============================================================

.PHONY: lint
lint: ## 執行 linter
	@which golangci-lint > /dev/null || (echo "請先安裝 golangci-lint" && exit 1)
	golangci-lint run ./...

.PHONY: fmt
fmt: ## 格式化程式碼
	$(GO) fmt ./...
	@which goimports > /dev/null && goimports -w . || true

.PHONY: vet
vet: ## 執行 go vet
	$(GO) vet ./...

.PHONY: tidy
tidy: ## 整理 go.mod
	$(GO) mod tidy

.PHONY: check
check: fmt vet lint test ## 執行所有檢查

# ============================================================
# 清理
# ============================================================

.PHONY: clean
clean: ## 清理建置產物
	rm -rf $(BINARY_DIR)
	rm -f coverage.out coverage.html

# ============================================================
# 輔助
# ============================================================

.PHONY: deps
deps: ## 下載依賴
	$(GO) mod download

.PHONY: help
help: ## 顯示說明
	@echo "VictoriaLogs MCP Server - Makefile 指令"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
