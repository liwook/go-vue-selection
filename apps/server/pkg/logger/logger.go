package logger

import (
	"os"

	"log/slog"

	"gopkg.in/natefinch/lumberjack.v2"
)

// Config 是 logger 包自有的配置结构，与上层 config 包解耦，
// 由组合根（main）负责将业务配置映射进来。
//
// 为什么不直接用 config 包的日志配置结构体当入参？
//  1. 依赖方向：logger 是底层可复用组件，不能 import 上层的 config 包，
//     否则会造成反向依赖甚至循环依赖（config/config 又可能被别的上层依赖）。
//  2. 稳定契约：自有 Config 作为 logger 对外的稳定接口，隔离内部实现细节。
//     即使将来 config 包的字段改名/增删，logger 也无需改动，只需组合根重做映射。
//  3. 可移植性：logger 包因此可以脱离本项目单独使用，不绑定任何业务配置形态。
//     结论：「能依赖 config」指的是上层依赖底层合法，但不代表该把别人的结构体
//     直接当成本组件的入参——入参用自有结构才是低耦合的做法。
type Config struct {
	Level      string // debug | info | warn | error（默认 info）
	Output     string // file | stdout | both（默认 file）
	Filename   string
	MaxSize    int
	MaxAge     int
	MaxBackups int
}

// levelVar 是并发安全的日志级别变量，支持运行期热更新（如 viper 监听配置变更）。
var levelVar = &slog.LevelVar{}

// Init 初始化Logger，接收 logger 自有的 Config，不依赖任何上层配置包。
func Init(cfg Config) (err error) {
	levelVar.Set(LevelFromString(cfg.Level))

	opts := &slog.HandlerOptions{
		Level:     levelVar,
		AddSource: true,
	}

	// 文件 handler（JSON 格式，按 lumberjack 切割）—— 裸机/VM 部署时使用
	fileHandler := slog.NewJSONHandler(&lumberjack.Logger{
		Filename:   cfg.Filename,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
	}, opts)
	// 终端 handler（text 格式）—— Docker/K8s 部署时由容器日志系统采集
	consoleHandler := slog.NewTextHandler(os.Stdout, opts)

	// 根据 output 配置决定日志去向：
	//   file   —— 仅写文件（默认，适合裸机/VM，由 lumberjack 切割）
	//   stdout —— 仅打 stdout（适合 Docker/K8s，docker logs 可见）
	//   both   —— 文件 + 终端（本地开发想两边都有）
	var handler slog.Handler
	switch cfg.Output {
	case "stdout":
		handler = consoleHandler
	case "both":
		handler = slog.NewMultiHandler(fileHandler, consoleHandler)
	default: // "file" 或为空
		handler = fileHandler
	}

	// 替换标准库全局的 default logger
	slog.SetDefault(slog.New(handler))

	return err
}

// SetLevel 运行期热更新日志级别（并发安全）。
func SetLevel(l slog.Level) {
	levelVar.Set(l)
}

func LevelFromString(s string) slog.Level {
	switch s {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default: // "info"、空串或未知值
		return slog.LevelInfo
	}
}
