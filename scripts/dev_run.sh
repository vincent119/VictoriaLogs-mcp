#!/bin/bash
set -e

# 顏色定義
GREEN='\033[0;32m'
NC='\033[0m'

echo -e "${GREEN}==> 啟動 VictoriaLogs MCP Server (開發模式)...${NC}"

# 檢查設定檔
if [ ! -f "configs/config.yaml" ]; then
    echo "警告: configs/config.yaml 不存在，使用 config.example.yaml 建立..."
    cp configs/config.example.yaml configs/config.yaml
fi

# 執行
go run cmd/vlmcp/main.go --config configs/config.yaml
