package handler

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/liwook/go-vue-selection/pkg/result"
)

type fileHandler struct {
	staticPath string
}

func NewFileHandler(staticPath string) *fileHandler {
	return &fileHandler{staticPath: staticPath}
}

// RegisterRoutes 注册文件上传路由（需在 JWT 中间件之后调用）
func (f *fileHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/file/upload", f.FileUpload)
}

// FileUpload 上传文件
// @Summary 上传文件接口
// @Description 上传文件
// @Tags 文件
// @Accept multipart/form-data
// @Produce application/json
// @Param file formData file true "文件"
// @Security ApiKeyAuth
// @Success 200 {object} result.ResponseData{data=string}
// @Router /api/v1/product/file/upload [post]
// @Failure 400 {object} result.ResponseData
// @Failure 500 {object} result.ResponseData
func (f *fileHandler) FileUpload(c *gin.Context) {
	// 限制请求体大小，避免超大文件（FormFile 会先整文件读入内存）导致 OOM
	const maxUploadSize = 10 << 20 // 10MB
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxUploadSize)

	file, err := c.FormFile("file")
	if err != nil {
		result.Error(c, result.CodeInvalidParam)
		return
	}

	// 仅允许图片类型，防止上传可执行文件等恶意内容
	ext := strings.ToLower(filepath.Ext(file.Filename))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".webp":
	default:
		slog.Warn("拒绝非图片类型上传", slog.String("filename", file.Filename))
		result.Error(c, result.CodeInvalidParam)
		return
	}

	// 取安全的纯文件名，防止路径穿越
	filename := filepath.Base(file.Filename)
	name := strings.TrimSuffix(filename, ext)
	savedName := fmt.Sprintf("%s_%d%s", name, time.Now().UnixNano(), ext)

	dateDir := time.Now().Format("20060102")
	// 直接用静态根目录拼接，避免 filepath.Base 丢失目录信息
	dist := filepath.Join(f.staticPath, "img", dateDir, savedName)
	if err = os.MkdirAll(filepath.Dir(dist), 0o755); err != nil {
		slog.Error("mkdir failed", slog.Any("error", err))
		result.Error(c, result.CodeServerBusy)
		return
	}
	// 上传文件到目标
	if err = c.SaveUploadedFile(file, dist); err != nil {
		slog.Error("save file failed", slog.Any("error", err))
		result.Error(c, result.CodeServerBusy)
		return
	}

	imgUrl := fmt.Sprintf("/api/img/%s/%s", dateDir, savedName)
	result.Success(c, imgUrl)
}
