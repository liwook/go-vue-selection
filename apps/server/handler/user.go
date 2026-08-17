package handler

import (
	"context"
	"github.com/liwook/go-vue-selection/pkg/errs"
	"github.com/liwook/go-vue-selection/pkg/translation"
	"github.com/liwook/go-vue-selection/types"
	"log/slog"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/liwook/go-vue-selection/pkg/result"
)

// UserService 在使用方（controller）定义接口，依赖倒置
type UserService interface {
	Login(context.Context, *types.ParamUserLogin) (string, error)
	SignUp(context.Context, *types.ParamUserSignUp) error
	GetUserInfo(context.Context, int64) (*types.ResponseUserInfo, error)
	GetUserList(context.Context, string, int64, int64) (*types.ResponseUserList, error)
	UpdateUser(context.Context, int64, *types.ParamUserUpdate) error
	DeleteUserByAdmin(context.Context, int64, int64) error
	DeleteSelfAccount(context.Context, int64) error
	ToAssign(context.Context, int64) (*types.ResponseToAssignRole, error)
	DoAssignRole(context.Context, *types.ParamDoAssignRole) error
	LockUser(context.Context, *types.ParamUserLock) error
}

type userHandler struct {
	userSvc UserService
}

func NewUserHandler(svc UserService) *userHandler {
	return &userHandler{userSvc: svc}
}

// RegisterPublicRoutes 注册无需 JWT 即可访问的路由（登录在 JWT 中间件之前调用）
func (u *userHandler) RegisterPublicRoutes(r *gin.RouterGroup) {
	r.POST("/index/login", u.Login)
}

// RegisterRoutes 注册用户管理相关路由（需在 JWT 中间件之后调用）
func (u *userHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/index/logout", u.Logout)
	r.GET("/user/info", u.GetInfo)
	r.POST("/user", u.SignUp)
	r.GET("/user", u.GetUser)
	r.GET("/user/:userId/role", u.ToAssign)
	r.POST("/user/:userId/role", u.DoAssignRole)
	r.PUT("/user", u.UpdateUser)
	r.DELETE("/user", u.DeleteSelf)
	r.DELETE("/user/:userId", u.DeleteUser)
	r.POST("/user/lock", u.LockUser)
}

// Login 处理登录请求的函数
// @Summary 用户登录接口
// @Description 处理用户登录请求
// @Tags 用户
// @Accept application/json
// @Produce application/json
// @Param object body types.ParamUserLogin true "用户登录参数"
// @Success 200 {object} result.ResponseData{data=string}
// @Router /api/v1/acl/index/login [post]
// @Failure 500 {object} result.ResponseData
func (u *userHandler) Login(c *gin.Context) {
	// 1. 获取参数及参数校验
	p := new(types.ParamUserLogin)
	if err := c.ShouldBindJSON(p); err != nil {
		// 请求参数有误，直接返回响应
		slog.Error("Login with invalid param", slog.Any("error", err))
		result.Error(c, result.CodeInvalidParam)
		return
	}
	// 2. 业务处理
	token, err := u.userSvc.Login(c.Request.Context(), p)
	if err != nil {
		slog.Error("Login failed", slog.String("username", p.Username), slog.Any("error", err))
		result.Error(c, errs.CodeOf(err))
		return
	}
	// 3. 返回响应
	result.Success(c, token)
}

// SignUp 处理新增用户接口
// @Summary 用户新增接口
// @Description 处理用户新增请求
// @Tags 用户
// @Accept application/json
// @Produce application/json
// @Param object body types.ParamUserSignUp true "用户注册信息"
// @Security ApiKeyAuth
// @Success 200 {object} result.ResponseData{data=string}
// @Router /api/v1/acl/user [post]
// @Failure 500 {object} result.ResponseData
func (u *userHandler) SignUp(c *gin.Context) {
	// 1. 参数校验
	p := new(types.ParamUserSignUp)
	if err := c.ShouldBindJSON(p); err != nil {
		slog.Error("SignUp with invalid param", slog.Any("error", err))
		// 判断 err 是不是 validator.ValidationErrors 类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			// 非 validator.validationErrors 类型错误直接返回
			result.Error(c, result.CodeInvalidParam)
			return
		}
		// validator.ValidationErrors类型错误则进行翻译
		errMap := translation.RemoveTopStruct(errs.Translate(translation.Trans))
		msgs := make([]string, 0, len(errMap))
		for _, v := range errMap {
			msgs = append(msgs, v)
		}
		result.ErrorWithMsg(c, result.CodeInvalidParam, strings.Join(msgs, "; "))
		return
	}

	// 2. 业务处理
	if err := u.userSvc.SignUp(c.Request.Context(), p); err != nil {
		slog.Error("SignUp failed", slog.Any("error", err))
		result.Error(c, errs.CodeOf(err))
		return
	}
	// 3. 返回响应
	result.Success(c, nil)
}

// GetInfo 获取用户信息处理函数
// @Summary 用户信息接口
// @Description 获取用户信息
// @Tags 用户
// @Accept application/json
// @Produce application/json
// @Security ApiKeyAuth
// @Success 200 {object} result.ResponseData{data=types.ResponseUserInfo}
// @Router /api/v1/acl/user/info [get]
// @Failure 500 {object} result.ResponseData
func (u *userHandler) GetInfo(c *gin.Context) {
	// 1. 获取当前用户ID
	userID, err := getCurrentUserID(c)
	if err != nil {
		result.Error(c, result.CodeNeedLogin)
		return
	}

	// 2. 获取用户信息
	userInfo, err := u.userSvc.GetUserInfo(c.Request.Context(), userID)
	if err != nil {
		result.Error(c, errs.CodeOf(err))
		return
	}

	result.Success(c, userInfo)
}

// Logout 处理登出函数
// @Summary 用户登出接口
// @Description 处理用户登出
// @Tags 用户
// @Accept application/json
// @Produce application/json
// @Security ApiKeyAuth
// @Success 200 {object} result.ResponseData{data=object}
// @Router /api/v1/acl/index/logout [post]
// @Failure 500 {object} result.ResponseData
func (u *userHandler) Logout(c *gin.Context) {
	// 1. 获取当前用户ID
	userID, err := getCurrentUserID(c)
	if err != nil {
		result.Error(c, result.CodeNeedLogin)
		return
	}

	// 2. 获取用户名并记录登出审计日志
	userInfo, err := u.userSvc.GetUserInfo(c.Request.Context(), userID)
	if err == nil && userInfo != nil {
		slog.Info("用户登出", slog.String("username", userInfo.Name), slog.Int64("user_id", userID))
	}

	result.Success(c, nil)
}

// GetUser 获取用户列表
// @Summary 获取用户分页列表
// @Description 获取用户列表
// @Tags 用户
// @Accept application/json
// @Produce application/json
// @Param page query int true "当前页码"
// @Param limit query int true "每页记录数"
// @Param username query string false "用户名"
// @Security ApiKeyAuth
// @Success 200 {object} result.ResponseData{data=types.ResponseUserList}
// @Router /api/v1/acl/user [get]
// @Failure 500 {object} result.ResponseData
func (u *userHandler) GetUser(c *gin.Context) {
	// 获取分页参数
	page, size := getPageInfo(c)
	username := c.Query("username")

	// 获取数据
	data, err := u.userSvc.GetUserList(c.Request.Context(), username, page, size)
	if err != nil {
		result.Error(c, errs.CodeOf(err))
		return
	}

	// 返回响应
	result.Success(c, data)
}

// UpdateUser 更新用户
// @Summary 更新用户接口
// @Description 处理更新用户请求
// @Tags 用户
// @Accept application/json
// @Produce application/json
// @Param object body types.ParamUserUpdate true "用户信息"
// @Security ApiKeyAuth
// @Success 200 {object} result.ResponseData{data=object}
// @Router /api/v1/acl/user [put]
// @Failure 500 {object} result.ResponseData
func (u *userHandler) UpdateUser(c *gin.Context) {
	p := new(types.ParamUserUpdate)
	if err := c.ShouldBindJSON(p); err != nil {
		slog.Error("UpdateUser with invalid param", slog.Any("error", err))
		result.Error(c, result.CodeInvalidParam)
		return
	}
	// user_id 由 JWT 令牌解析得到，避免前端篡改
	userID, err := getCurrentUserID(c)
	if err != nil {
		slog.Error("UpdateUser without login", slog.Any("error", err))
		result.Error(c, result.CodeNeedLogin)
		return
	}
	err = u.userSvc.UpdateUser(c.Request.Context(), userID, p)
	if err != nil {
		result.Error(c, errs.CodeOf(err))
		return
	}
	result.Success(c, nil)
}

// DeleteUser 删除用户
// @Summary 删除用户接口
// @Description 处理删除用户请求
// @Tags 用户
// @Accept application/json
// @Produce application/json
// @Param userId path string true "用户 ID"
// @Security ApiKeyAuth
// @Success 200 {object} result.ResponseData{data=object}
// @Router /api/v1/acl/user/{userId} [delete]
// @Failure 500 {object} result.ResponseData
func (u *userHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("userId")
	targetUserID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		slog.Error("DeleteUser with invalid param", slog.Any("error", err))
		result.Error(c, result.CodeInvalidParam)
		return
	}

	// 操作者 ID 从 JWT 令牌解析得到
	operatorID, err := getCurrentUserID(c)
	if err != nil {
		slog.Error("DeleteUser without login", slog.Any("error", err))
		result.Error(c, result.CodeNeedLogin)
		return
	}
	// 禁止删除自己，避免管理员把自己锁死
	if operatorID == targetUserID {
		slog.Error("DeleteUser reject self-delete", slog.Int64("user_id", targetUserID))
		result.Error(c, result.CodeNoPermission)
		return
	}

	err = u.userSvc.DeleteUserByAdmin(c.Request.Context(), operatorID, targetUserID)
	if err != nil {
		result.Error(c, errs.CodeOf(err))
		return
	}
	result.Success(c, nil)
}

// DeleteSelf 注销当前登录用户（删除自己）
// @Summary 注销当前用户接口
// @Description 删除当前登录用户自己，用户 ID 从 JWT 令牌解析
// @Tags 用户
// @Accept application/json
// @Produce application/json
// @Security ApiKeyAuth
// @Success 200 {object} result.ResponseData{data=object}
// @Router /api/v1/acl/user [delete]
// @Failure 500 {object} result.ResponseData
func (u *userHandler) DeleteSelf(c *gin.Context) {
	// user_id 由 JWT 令牌解析得到，避免前端篡改
	userID, err := getCurrentUserID(c)
	if err != nil {
		slog.Error("DeleteSelf without login", slog.Any("error", err))
		result.Error(c, result.CodeNeedLogin)
		return
	}

	err = u.userSvc.DeleteSelfAccount(c.Request.Context(), userID)
	if err != nil {
		result.Error(c, errs.CodeOf(err))
		return
	}
	result.Success(c, nil)
}

// ToAssign 获取用户角色分配接口
// @Summary 用户角色分配接口
// @Description 获取用户角色分配息
// @Tags 用户
// @Accept application/json
// @Produce application/json
// @Param userId path string true "用户 ID"
// @Security ApiKeyAuth
// @Success 200 {object} result.ResponseData{data=types.ResponseToAssignRole}
// @Router /api/v1/acl/user/{userId}/role [get]
// @Failure 500 {object} result.ResponseData
func (u *userHandler) ToAssign(c *gin.Context) {
	idStr := c.Param("userId")
	userId, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		slog.Error("ToAssign with invalid param", slog.Any("error", err))
		result.Error(c, result.CodeInvalidParam)
		return
	}
	data, err := u.userSvc.ToAssign(c.Request.Context(), userId)
	if err != nil {
		result.Error(c, errs.CodeOf(err))
		return
	}
	result.Success(c, data)
}

// DoAssignRole 为用户分配角色
// @Summary 为用户分配角色
// @Description 为用户分配角色
// @Tags 用户
// @Accept application/json
// @Produce application/json
// @Param userId path string true "用户 ID"
// @Param object body types.ParamDoAssignRole true "角色分配信息"
// @Security ApiKeyAuth
// @Success 200 {object} result.ResponseData{data=object}
// @Router /api/v1/acl/user/{userId}/role [post]
// @Failure 500 {object} result.ResponseData
func (u *userHandler) DoAssignRole(c *gin.Context) {
	p := new(types.ParamDoAssignRole)
	if err := c.ShouldBindJSON(p); err != nil {
		slog.Error("DoAssignRole with invalid param", slog.Any("error", err))
		result.Error(c, result.CodeInvalidParam)
		return
	}
	userId, err := strconv.ParseInt(c.Param("userId"), 10, 64)
	if err != nil {
		slog.Error("DoAssignRole with invalid param", slog.Any("error", err))
		result.Error(c, result.CodeInvalidParam)
		return
	}
	p.UserID = strconv.FormatInt(userId, 10)
	err = u.userSvc.DoAssignRole(c.Request.Context(), p)
	if err != nil {
		result.Error(c, errs.CodeOf(err))
		return
	}
	result.Success(c, nil)
}

// LockUser 锁定/解锁用户
// @Summary 锁定/解锁用户接口
// @Description 根据 status 锁定（false）或解锁（true）用户账号
// @Tags 用户
// @Accept application/json
// @Produce application/json
// @Param object body types.ParamUserLock true "锁定/解锁参数"
// @Security ApiKeyAuth
// @Success 200 {object} result.ResponseData{data=object}
// @Router /api/v1/acl/user/lock [post]
// @Failure 500 {object} result.ResponseData
func (u *userHandler) LockUser(c *gin.Context) {
	p := new(types.ParamUserLock)
	if err := c.ShouldBindJSON(p); err != nil {
		slog.Error("LockUser with invalid param", slog.Any("error", err))
		result.Error(c, result.CodeInvalidParam)
		return
	}
	if err := u.userSvc.LockUser(c.Request.Context(), p); err != nil {
		slog.Error("LockUser failed", slog.Any("error", err))
		result.Error(c, errs.CodeOf(err))
		return
	}
	result.Success(c, nil)
}
