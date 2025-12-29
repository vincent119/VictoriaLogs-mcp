# 發布流程 (Release Process)

本專案使用 `make` 與 Shell Scripts 進行自動化發布建置。

## 版本策略

遵循 [Semantic Versioning 2.0.0](https://semver.org/)。

- **Major**: 不相容的 API 變更。
- **Minor**: 向下相容的功能新增。
- **Patch**: 向下相容的問題修正。

## 建置流程

### 1. 本地測試

在發布前，務必執行完整測試：

```bash
make lint
make test
```

### 2. 執行發布建置

使用 `scripts/release_build.sh` 腳本進行跨平台編譯：

```bash
./scripts/release_build.sh v1.0.0
```

此腳本會：
1. 檢查 git 狀態是否乾淨。
2. 執行測試。
3. 編譯 macOS (arm64/amd64) 與 Linux (amd64/arm64) 版本。
4. 將二進制檔案輸出至 `bin/release/`。
5. 生成 SHA256 checksums。

### 3. 發布 Artifacts

產生的檔案結構：
```
bin/release/
├── vlmcp-v1.0.0-darwin-arm64
├── vlmcp-v1.0.0-linux-amd64
└── checksums.txt
```

將這些檔案上傳至 GitHub Releases 頁面。
