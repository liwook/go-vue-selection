package errs

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/liwook/go-vue-selection/pkg/result"
)

// BizError 携带业务错误码的错误，由 service 层在 repo 裸 error 之上包装。
// 它实现了 error 与 Unwrap 接口，便于 errors.As / errors.Is 提取与透传。
type BizError struct {
	Code result.ResCode
	Err  error
}

func (e *BizError) Error() string { return e.Err.Error() }

func (e *BizError) Unwrap() error { return e.Err }

// 错误处理约定（service 在 repo 边界处"打印 + 返回"）：
//   - repo 层：
//     · 单 DB 操作的函数：直接返回裸 error，不包装（步骤语义由 service 层体现）。
//     · 多 DB 操作（或 DB + bcrypt/计数等混合操作）的函数：每个 return err 用
//     fmt.Errorf("整体动作(具体步骤): %w", err) 包一层纯文本语义，便于精确定位
//     到底是哪一步挂了；repo 不决定业务码、不打印。
//     · 不打印（数据库驱动层也不打印）。
//   - service 层（repo 裸 error 之上）：作为错误真正发生的边界，
//     必须「slog 打印一次（带业务上下文）+ errs.Wrap 返回」。
//     打印是可观测性，返回是控制流，两者职责不同、不冲突。
//   - handler 层：只调用 result.Error(c, errs.CodeOf(err)) 映射错误码，不打印，
//     避免双重打印。
//
// 整条链路只有 service 这一层打印；若 service 也只返回不打印，错误将无人记录。
//
// Wrap 将一个底层 error 包装为带业务码的错误。err 为 nil 时返回 nil。
func Wrap(code result.ResCode, err error) error {
	if err == nil {
		return nil
	}
	return &BizError{Code: code, Err: err}
}

// CodeOf 从 error 中提取业务错误码。
// - nil 返回成功码；
// - 若是 *BizError 返回其携带的码；
// - 否则（裸 error / 哨兵 error）回退到系统繁忙码。
func CodeOf(err error) result.ResCode {
	if err == nil {
		return result.CodeSuccess
	}

	be, ok := errors.AsType[*BizError](err)
	if ok {
		return be.Code
	}

	return result.CodeServerBusy
}

// PostgreSQL SQLSTATE 错误码（常用子集）
// 见 https://www.postgresql.org/docs/current/errcodes-appendix.html
const (
	sqlStateForeignKeyViolation = "23503" // foreign_key_violation 违反外键约束
)

// IsForeignKeyViolation 判断错误链中是否包含 PostgreSQL 外键约束冲突
// （SQLSTATE 23503），即"被其他表引用的资源不允许删除"（ON DELETE RESTRICT 一侧）。
// errors.AsType 会穿透 errs.Wrap 包装直达驱动层 *pgconn.PgError。
func IsForeignKeyViolation(err error) bool {
	if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
		return pgErr.Code == sqlStateForeignKeyViolation
	}
	return false
}
