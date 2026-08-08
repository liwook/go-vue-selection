//go:build tools
// +build tools

// 该文件仅用于锁定开发工具依赖（如 swag），不参与正常编译。
package tools

import (
	_ "github.com/swaggo/swag/cmd/swag"
)
