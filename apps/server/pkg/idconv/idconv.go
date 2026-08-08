package idconv

import (
	"strconv"
)

// ToStr 将 int64 类型的 ID 转换为字符串（用于响应体对外暴露）。
func ToStr(id int64) string {
	return strconv.FormatInt(id, 10)
}

// ToStrPtr 将 *int64 类型的可选 ID 转换为字符串；nil 时返回空串
// （区别于 ToStr(derefInt64(nil)) 会得到 "0"，可选字段无值时语义上应为空串）。
func ToStrPtr(id *int64) string {
	if id == nil {
		return ""
	}
	return strconv.FormatInt(*id, 10)
}

// ToInt64 将字符串 ID 转换为 int64，解析失败返回 error。
func ToInt64(s string) (int64, error) {
	if s == "" {
		return 0, nil
	}
	return strconv.ParseInt(s, 10, 64)
}

// ToInt64Safe 将字符串 ID 转换为 int64，空串或非法字符一律返回 0（用于可选/容忍场景）。
func ToInt64Safe(s string) int64 {
	v, err := ToInt64(s)
	if err != nil {
		return 0
	}
	return v
}
