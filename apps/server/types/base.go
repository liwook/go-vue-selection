package types

import (
	"time"
)

// BaseModel 通用基类
// CreateTime/UpdateTime 使用 time.Time，JSON 序列化默认输出 RFC3339（带时区偏移）
// 例如 "2026-08-03T15:04:05+08:00"，前端用 new Date() 可无歧义解析
// 零值时由 omitzero 省略（Go 1.24+），避免返回 "0001-01-01T00:00:00Z" 垃圾时间
type BaseModel struct {
	CreateTime time.Time `json:"createTime,omitzero"`
	UpdateTime time.Time `json:"updateTime,omitzero"`
}
