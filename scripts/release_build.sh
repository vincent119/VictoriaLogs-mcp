#!/bin/bash
set -e

VERSION=$1

if [ -z "$VERSION" ]; then
    echo "Usage: $0 <version>"
    echo "Example: $0 v1.0.0"
    exit 1
fi

OUTPUT_DIR="bin/release"
mkdir -p $OUTPUT_DIR

echo "==> 開始發布建置: $VERSION"

# 1. 執行測試
echo "==> 執行測試..."
go test ./...

# 2. 定義目標平台
PLATFORMS=("darwin/arm64" "darwin/amd64" "linux/amd64" "linux/arm64")

for PLATFORM in "${PLATFORMS[@]}"; do
    GOOS=${PLATFORM%/*}
    GOARCH=${PLATFORM#*/}
    OUTPUT_NAME="vlmcp-${VERSION}-${GOOS}-${GOARCH}"

    echo "==> 編譯 $GOOS/$GOARCH..."
    env GOOS=$GOOS GOARCH=$GOARCH go build -ldflags "-s -w -X github.com/vincent119/victorialogs-mcp/pkg/version.Version=${VERSION}" -o "${OUTPUT_DIR}/${OUTPUT_NAME}" ./cmd/vlmcp
done

# 3. 產生 Checksums
echo "==> 產生 Checksums..."
cd $OUTPUT_DIR
shasum -a 256 * > checksums.txt

echo "==> 發布建置完成！檔案位於 $OUTPUT_DIR"
ls -lh
