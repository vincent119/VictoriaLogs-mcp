// Package version 提供應用程式版本資訊
package version

import (
	"fmt"
	"runtime"
)

// 版本資訊變數，透過 ldflags 在編譯時注入
var (
	// Version 版本號（透過 git tag 或手動設定）
	Version = "dev"

	// GitCommit Git commit hash
	GitCommit = "unknown"

	// BuildTime 建置時間
	BuildTime = "unknown"
)

// Info 版本資訊結構
type Info struct {
	Version   string `json:"version"`
	GitCommit string `json:"git_commit"`
	BuildTime string `json:"build_time"`
	GoVersion string `json:"go_version"`
	OS        string `json:"os"`
	Arch      string `json:"arch"`
}

// Get 取得版本資訊
func Get() Info {
	return Info{
		Version:   Version,
		GitCommit: GitCommit,
		BuildTime: BuildTime,
		GoVersion: runtime.Version(),
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
	}
}

// String 回傳版本字串
func (i Info) String() string {
	return fmt.Sprintf(
		"Version: %s\nGit Commit: %s\nBuild Time: %s\nGo Version: %s\nOS/Arch: %s/%s",
		i.Version, i.GitCommit, i.BuildTime, i.GoVersion, i.OS, i.Arch,
	)
}

// Short 回傳簡短版本字串
func Short() string {
	return Version
}
