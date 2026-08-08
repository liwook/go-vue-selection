package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
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

	// 7. 启动服务（含优雅关闭）
	run(r, cfg.Port)
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
