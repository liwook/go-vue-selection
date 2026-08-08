package seed

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/liwook/go-vue-selection/dal/model"
	"github.com/liwook/go-vue-selection/dal/query"
	"github.com/liwook/go-vue-selection/pkg/snowflake"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Run 初始化种子数据：仅在 users 表为空时执行，保证幂等（容器重启/重复启动不会重复插入）。
// 带雪花ID的业务数据（users.user_id、user_role）由代码生成并写入；
// 其余字典数据（role、menu、category 等固定ID）仍由 init-sql 负责。
func Run(db *gorm.DB) error {
	q := query.Use(db)
	ctx := context.Background()

	// 幂等判断：users 表已有数据则跳过
	cnt, err := q.User.WithContext(ctx).Count()
	if err != nil {
		return fmt.Errorf("seed: count users failed: %w", err)
	}
	if cnt > 0 {
		slog.Info("seed skipped: users table already has data", "count", cnt)
		return nil
	}

	// 密码统一使用 bcrypt 哈希（原 SQL 里是明文 md5，不符合安全规范）
	hash, err := bcrypt.GenerateFromPassword([]byte("111111"), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("seed: gen password hash failed: %w", err)
	}

	// user_id 走和线上一致的雪花算法（不写死，避免与业务生成的 ID 冲突）
	adminID := snowflake.GenID()
	userID := snowflake.GenID()

	avatarAdmin := "/api/static/img/avatar/admin.gif"
	avatarUser := "/api/static/img/avatar/user.png"
	phoneAdmin := "13800138000"
	phoneUser := "13800138001"

	// 插入默认用户（users 表以 user_id 为主键，只填业务字段）
	if err := q.User.WithContext(ctx).Create(
		&model.User{
			UserID:   adminID,
			Username: "admin",
			Password: string(hash),
			Name:     new("管理员"),
			Phone:    &phoneAdmin,
			Avatar:   &avatarAdmin,
		},
		&model.User{
			UserID:   userID,
			Username: "user",
			Password: string(hash),
			Name:     new("李四"),
			Phone:    &phoneUser,
			Avatar:   &avatarUser,
		},
	); err != nil {
		return fmt.Errorf("seed: insert users failed: %w", err)
	}

	// 插入用户-角色关联（role_id 来自 init-sql 的固定字典：1=admin, 2=user）
	if err := q.UserRole.WithContext(ctx).Create(
		&model.UserRole{UserID: adminID, RoleID: 1},
		&model.UserRole{UserID: userID, RoleID: 2},
	); err != nil {
		return fmt.Errorf("seed: insert user_role failed: %w", err)
	}

	slog.Info("seed done: inserted default users and user_role")
	return nil
}

// Init 按序执行全部种子数据初始化：先写入基础用户/角色，再写入商城演示数据。
// 各子函数内部已各自实现幂等（按对应表是否为空判断），重复启动不会重复插入。
func Init(db *gorm.DB) error {
	if err := Run(db); err != nil {
		return fmt.Errorf("seed init: %w", err)
	}
	if err := RunDemo(db); err != nil {
		return fmt.Errorf("seed init: %w", err)
	}
	return nil
}
