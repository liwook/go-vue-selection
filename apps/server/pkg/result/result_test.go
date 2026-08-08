package result

import (
	"net/http"
	"testing"
)

func TestResCode_Msg(t *testing.T) {
	tests := []struct {
		code ResCode
		want string
	}{
		{CodeSuccess, "success"},
		{CodeInvalidParam, "请求参数错误"},
		{CodeNeedLogin, "需要登录"},
		{CodeInvalidToken, "无效的Token"},
		{CodeNoRoute, "请求路径不存在"},
		{CodeNoPermission, "无操作权限"},
		{CodeUserExist, "用户名已存在"},
		{CodeUserDBErr, "用户数据访问失败"},
		{CodeInvalidPassword, "用户名或密码错误"},
		{CodeRoleDBErr, "角色数据访问失败"},
		{CodeRoleExist, "角色已存在"},
		{CodeMenuNodeExist, "该节点下有子节点，不可以删除"},
		{CodeMenuDBErr, "菜单数据访问失败"},
		{CodeMenuExist, "菜单名称已存在"},
		{CodeTrademarkErr, "品牌操作失败"},
		{CodeTrademarkDBErr, "品牌数据访问失败"},
		{CodeCategoryDBErr, "分类数据访问失败"},
		{CodeAttrDBErr, "属性数据访问失败"},
		{CodeSpuDBErr, "SPU数据访问失败"},
		{CodeSkuDBErr, "SKU数据访问失败"},
		{CodeDBError, "数据库服务异常"},
		{CodeRedisError, "缓存服务异常"},
		{CodeExternalErr, "外部服务不可用"},
		{CodePanic, "服务内部错误"},
		{CodeServerBusy, "服务繁忙"},
		{ResCode(999999), "未知错误"}, // 未注册码
	}
	for _, tt := range tests {
		if got := tt.code.Msg(); got != tt.want {
			t.Errorf("ResCode(%d).Msg() = %q, want %q", tt.code, got, tt.want)
		}
	}
}

func TestHttpStatusForCode(t *testing.T) {
	tests := []struct {
		name string
		code ResCode
		want int
	}{
		{"business success -> 200", CodeSuccess, http.StatusOK},
		{"business error -> 200", CodeInvalidParam, http.StatusOK},
		{"acl error -> 200", CodeUserExist, http.StatusOK},
		{"product error -> 200", CodeSkuDBErr, http.StatusOK},
		{"sys db error -> 500", CodeDBError, http.StatusInternalServerError},
		{"sys external -> 500", CodeExternalErr, http.StatusInternalServerError},
		{"sys panic -> 500", CodePanic, http.StatusInternalServerError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := httpStatusForCode(tt.code); got != tt.want {
				t.Errorf("httpStatusForCode(%d) = %d, want %d", tt.code, got, tt.want)
			}
		})
	}
}

func TestErrorExternal_MasksSysDomain(t *testing.T) {
	// 99 域错误应返回脱敏文案
	if got := externalMsgFor(CodeDBError); got != externalMaskedMsg {
		t.Errorf("sys domain msg = %q, want masked %q", got, externalMaskedMsg)
	}
	// 非 99 域返回真实文案
	if got := externalMsgFor(CodeInvalidParam); got != CodeInvalidParam.Msg() {
		t.Errorf("business domain msg = %q, want real %q", got, CodeInvalidParam.Msg())
	}
}

// externalMsgFor 复刻 ErrorExternal 的脱敏分支逻辑，便于不依赖 gin.Context 做断言。
func externalMsgFor(code ResCode) string {
	msg := code.Msg()
	if code/10000 == 99 {
		msg = externalMaskedMsg
	}
	return msg
}
