//go:build integration

package tests

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/liwook/go-vue-selection/config"
	"github.com/liwook/go-vue-selection/pkg/snowflake"
	"github.com/liwook/go-vue-selection/repository"
	"github.com/liwook/go-vue-selection/router"
	"github.com/liwook/go-vue-selection/seed"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const testDBName = "vue_admin_test"

// TestMain 是集成测试的入口：自动建库 → 建表(init.sql) → 灌种子数据(seed) →
// 进程内用 httptest 拉起真实 router，供所有用例复用。
// 完全不触碰真实业务库 vue_admin，隔离在 vue_admin_test 里。
func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	// 1. 加载配置（与线上一致，连接信息来自 etc/config.yaml）
	cfg, err := config.Init("../etc/config.yaml")
	if err != nil {
		slog.Error("load config failed", "error", err)
		os.Exit(1)
	}

	// 2. 初始化雪花 ID（seed 依赖）
	if err := snowflake.Init(cfg.SnowFlake.MachineID); err != nil {
		slog.Error("init snowflake failed", "error", err)
		os.Exit(1)
	}

	// 3. 自动创建测试库 vue_admin_test（已存在则跳过）
	if err := ensureTestDatabase(cfg.Postgres); err != nil {
		slog.Error("ensure test database failed", "error", err)
		os.Exit(1)
	}

	// 4. 复制一份配置，但把 dbname 换成测试库
	testPg := cfg.Postgres
	testPg.DbName = testDBName

	// 5. 连接测试库，自动跑 init.sql 建表 + seed 灌数据
	db, err := repository.Init(&testPg)
	if err != nil {
		slog.Error("connect test database failed", "error", err)
		os.Exit(1)
	}
	defer repository.Close(db)

	if err := setupSchema(db); err != nil {
		slog.Error("setup schema failed", "error", err)
		os.Exit(1)
	}
	if err := seed.Init(db); err != nil {
		slog.Error("seed failed", "error", err)
		os.Exit(1)
	}

	// 6. 进程内拉起真实 router
	r := router.Setup(cfg, db)

	// 7. 初始化 API 客户端（含登录拿 Token）
	if err := setupAPIClient(r, &testPg); err != nil {
		slog.Error("setup api client failed", "error", err)
		os.Exit(1)
	}

	os.Exit(m.Run())
}

// ensureTestDatabase 连到 postgres 默认库，检查 vue_admin_test 是否存在，不存在则 CREATE DATABASE。
func ensureTestDatabase(pg config.PostgresConfig) error {
	adminDSN := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=postgres sslmode=disable",
		pg.Host, pg.Port, pg.User, pg.Password)
	adminDB, err := gorm.Open(postgres.Open(adminDSN), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("open postgres admin db: %w", err)
	}
	sqlDB, err := adminDB.DB()
	if err != nil {
		return fmt.Errorf("get sql.DB: %w", err)
	}
	defer sqlDB.Close()

	var exists bool
	checkSQL := `SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)`
	if err := adminDB.Raw(checkSQL, testDBName).Scan(&exists).Error; err != nil {
		return fmt.Errorf("check database exists: %w", err)
	}
	if exists {
		slog.Info("test database already exists, skip create", "db", testDBName)
		return nil
	}

	if err := adminDB.Exec(fmt.Sprintf("CREATE DATABASE %s", testDBName)).Error; err != nil {
		return fmt.Errorf("create database %s: %w", testDBName, err)
	}
	slog.Info("test database created", "db", testDBName)
	return nil
}

// setupSchema 读取 init-sql/init.sql 并按语句拆分执行，完成建表/触发器/外键/字典数据。
// 注意保护 $$ 包裹的 plpgsql 函数体内的分号，避免被错误拆分。
func setupSchema(db *gorm.DB) error {
	sqlPath := filepath.Join("..", "init-sql", "init.sql")
	content, err := os.ReadFile(sqlPath)
	if err != nil {
		return fmt.Errorf("read init.sql: %w", err)
	}

	for _, stmt := range splitSQL(string(content)) {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		if err := db.Exec(stmt).Error; err != nil {
			return fmt.Errorf("exec init.sql statement failed: %w\nstmt: %s", err, stmt)
		}
	}
	slog.Info("init.sql executed, schema ready")
	return nil
}

// splitSQL 按分号拆分 SQL，但忽略 $$ ... $$ 块内的分号（plpgsql 函数体）。
func splitSQL(s string) []string {
	var stmts []string
	var buf strings.Builder
	inDollar := false
	for i := 0; i < len(s); i++ {
		c := s[i]
		if i+1 < len(s) && c == '$' && s[i+1] == '$' {
			inDollar = !inDollar
			buf.WriteByte(c)
			buf.WriteByte(s[i+1])
			i++
			continue
		}
		if c == ';' && !inDollar {
			stmts = append(stmts, buf.String())
			buf.Reset()
			continue
		}
		buf.WriteByte(c)
	}
	if strings.TrimSpace(buf.String()) != "" {
		stmts = append(stmts, buf.String())
	}
	return stmts
}
