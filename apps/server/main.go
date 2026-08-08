package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/liwook/go-vue-selection/config"
	"github.com/liwook/go-vue-selection/pkg/logger"
	"github.com/liwook/go-vue-selection/pkg/snowflake"
	"github.com/liwook/go-vue-selection/pkg/translation"
	"github.com/liwook/go-vue-selection/repository"
	"github.com/liwook/go-vue-selection/router"
	"github.com/liwook/go-vue-selection/seed"

	"log/slog"
)

// version 为程序版本号，构建时可通过 -ldflags "-X main.version=v1.2.3" 注入，
// 未注入时回退到此处默认值。
var version = "dev"

// @title vue_admin 项目接口文档
// @version 1.0
// @description 硅谷甄选项目后端

// @host 127.0.0.1:9000
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description 认证 Token，格式：Bearer <token>
/*
主函数，负责初始化应用程序的各个组件并启动HTTP服务器
包括配置加载、日志初始化、数据库连接、雪花ID生成器、校验器翻译、
种子数据初始化、路由注册以及服务的优雅关闭等功能
*/
func main() {
	// 0. 解析命令行参数（含 -v/--version）
	configPath, versionOnly := parseArgs()
	if versionOnly {
		fmt.Printf("vue_admin %s\n", version)
		return
	}

	// 1. 加载配置
	cfg, err := config.Init(configPath)
	if err != nil {
		fmt.Printf("init settings failed, err:%v\n", err)
		return
	}

	// 2. 初始化日志（含日志级别热更新钩子）
	if err := setupLogger(cfg); err != nil {
		return
	}
	slog.Debug("logger init success...")

	// 3. 初始化PostgreSQL连接
	db, err := repository.Init(&cfg.Postgres)
	if err != nil {
		fmt.Printf("init postgres failed, err:%v\n", err)
		return
	}
	defer repository.Close(db)

	// 4. 初始化雪花ID
	if err := snowflake.Init(cfg.SnowFlake.MachineID); err != nil {
		fmt.Printf("init failed, err:%v \n", err)
		return
	}

	// 5. 初始化gin框架内置校验器的翻译器
	if err := translation.InitTrans("zh"); err != nil {
		fmt.Printf("init validator trans failed, err:%v \n", err)
		return
	}

	// 5.5 初始化种子数据（基础用户/角色 + 商城演示数据，各自幂等）
	if err := seed.Init(db); err != nil {
		fmt.Printf("seed failed, err:%v\n", err)
		return
	}

	// 6. 注册路由
	r := router.Setup(cfg, db)

	// 6.5 启动独立的 pprof 调试 server（常驻、仅绑定 localhost，不污染主业务端口）
	startPprofServer(cfg)

	// 7. 启动服务（含优雅关闭）
	run(r, cfg.Port)
}

// startPprofServer 按配置启动一个独立的 pprof 调试 HTTP server。
// 采用独立 ServeMux + 绑定 127.0.0.1 的方案，确保：
//  1. 主业务端口（如 9000）不暴露 /debug/pprof，避免信息泄露与 DoS 风险；
//  2. 仅本机/跳板机可达，外网不可达；
//  3. 常驻但几乎零开销——只有主动发 /debug/pprof/profile 或 /trace 才会真正采样。
func startPprofServer(cfg *config.AppConfig) {
	if !cfg.Pprof.Enabled {
		slog.Debug("pprof disabled, skip starting debug server")
		return
	}
	port := cfg.Pprof.Port
	if port == 0 {
		port = 6060
	}
	// 独立 mux，避免污染 DefaultServeMux 与 gin 主路由。
	// 注：这里没有用 `import _ "net/http/pprof"` 的简写方式，原因是：
	//   - 该方式会把路由注册到全局的 http.DefaultServeMux；
	//   - 而本项目使用 gin，主 server 走的是 gin 自己的 router，并不接 DefaultServeMux，
	//     若依赖 DefaultServeMux 需另起一个 handler 为 nil 的 server 才能生效；
	//   - 显式 NewServeMux() 可将 debug 路由严格隔离在独立 mux 中，
	//     既不与任何全局默认 mux 纠缠，也保证主业务端口（port）绝不暴露 /debug/pprof。
	mux := http.NewServeMux()
	mux.HandleFunc("/debug/pprof/", pprof.Index)          // 概览页，列出所有 profile 类型
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline) // 启动命令行参数
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile) // CPU profile（默认采未来 30s）
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)   // 地址转符号
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)     // 执行跟踪

	addr := fmt.Sprintf("127.0.0.1:%d", port)
	srv := &http.Server{Addr: addr, Handler: mux}
	go func() {
		slog.Info("pprof debug server started", "addr", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("pprof server start failed", "error", err)
		}
	}()
}

// setupLogger 初始化日志组件并注册日志级别热更新钩子。
// 把这部分接线逻辑从 main 抽出，避免 main 承载过多细节；
// 钩子放在组合根而非 logger 包内，可避免 logger 反向依赖 config。
func setupLogger(cfg *config.AppConfig) error {
	if err := logger.Init(logger.Config{
		Level:      cfg.Log.Level,
		Output:     cfg.Log.Output,
		Filename:   cfg.Log.Filename,
		MaxSize:    cfg.Log.MaxSize,
		MaxAge:     cfg.Log.MaxAge,
		MaxBackups: cfg.Log.MaxBackups,
	}); err != nil {
		fmt.Printf("init logger failed, err:%v\n", err)
		return err
	}

	config.OnChange(func(c *config.AppConfig) {
		logger.SetLevel(logger.LevelFromString(c.Log.Level))
	})
	return nil
}

// parseArgs 解析命令行参数，返回配置文件路径与是否仅查看版本。
func parseArgs() (configPath string, versionOnly bool) {
	var showVersion bool
	flag.StringVar(&configPath, "f", "./etc/config.yaml", "配置文件路径")
	flag.BoolVar(&showVersion, "v", false, "查看版本信息")
	flag.BoolVar(&showVersion, "version", false, "查看版本信息")
	flag.Parse()
	return configPath, showVersion
}

// run 启动 HTTP 服务并在收到中断信号时优雅关闭。
func run(r http.Handler, port int) {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server start failed", "error", err)
			os.Exit(1)
		}
	}()

	// 等待中断信号来优雅地关闭服务器，为关闭服务器操作设置一个5秒的超时
	// kill 默认发送 syscall.SIGTERM；Ctrl+C 发送 syscall.SIGINT；SIGKILL 不可捕获故不监听
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop() // 退出前停止监听信号，释放资源
	<-ctx.Done() // 阻塞在此，收到上述信号后 ctx 被取消才会往下执行
	slog.Info("Shutdown Server ...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 5秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过5秒就超时退出
	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("Server Shutdown: ", "error", err)
		os.Exit(1)
	}

	slog.Info("Server exiting")
}
