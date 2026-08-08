package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"strconv"

	"github.com/liwook/go-vue-selection/pkg/middlewares"
)

var ErrorUserNotLogin = errors.New("用户未登录")

// getCurrentUserID 获取当前用户ID
func getCurrentUserID(c *gin.Context) (userID int64, err error) {
	uid, ok := c.Get(middlewares.CtxUserIDKey)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	userID, ok = uid.(int64)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	return
}

func getPageInfo(c *gin.Context) (int64, int64) {
	pageStr := c.DefaultQuery("page", "1")
	sizeStr := c.DefaultQuery("limit", "10")

	var (
		page  int64
		limit int64
		err   error
	)
	page, err = strconv.ParseInt(pageStr, 10, 64)
	if err != nil {
		page = 1
	}
	limit, err = strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		limit = 10
	}
	return page, limit
}

func getAllCategoryID(c *gin.Context) (int64, int64, int64, error) {
	c1IdStr := c.Param("category1Id")
	c2IdStr := c.Param("category2Id")
	c3IdStr := c.Param("category3Id")

	var (
		c1Id int64
		c2Id int64
		c3Id int64
		err  error
	)
	c1Id, err = strconv.ParseInt(c1IdStr, 10, 64)
	if err != nil {
		return 0, 0, 0, err
	}
	c2Id, err = strconv.ParseInt(c2IdStr, 10, 64)
	if err != nil {
		return 0, 0, 0, err
	}
	c3Id, err = strconv.ParseInt(c3IdStr, 10, 64)
	if err != nil {
		return 0, 0, 0, err
	}
	return c1Id, c2Id, c3Id, err
}
