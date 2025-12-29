// VictoriaLogs MCP Server
// 提供 MCP tools 給 AI/Agent 查詢 VictoriaLogs
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/vincent119/victorialogs-mcp/internal/app"
	"github.com/vincent119/victorialogs-mcp/internal/config"
	"github.com/vincent119/victorialogs-mcp/pkg/version"
)

var (
	configPath  = flag.String("config", "", "config path")
	showVersion = flag.Bool("version", false, "version")
)

func main() {
	flag.Parse()

	// 顯示版本資訊
	if *showVersion {
		fmt.Println(version.Get().String())
		os.Exit(0)
	}

	// 載入設定
	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load config failed: %v\n", err)
		os.Exit(1)
	}

	// 建立應用程式
	application, err := app.New(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "init app failed: %v\n", err)
		os.Exit(1)
	}

	// 執行應用程式（含優雅關閉）
	if err := app.RunWithGracefulShutdown(application); err != nil {
		fmt.Fprintf(os.Stderr, "run app failed: %v\n", err)
		os.Exit(1)
	}
}
