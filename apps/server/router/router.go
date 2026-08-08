package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"

	"github.com/liwook/go-vue-selection/config"
	"github.com/liwook/go-vue-selection/handler"
	"github.com/liwook/go-vue-selection/pkg/jwt"
	"github.com/liwook/go-vue-selection/pkg/middlewares"
	"github.com/liwook/go-vue-selection/pkg/result"
	"github.com/liwook/go-vue-selection/repository"
	"github.com/liwook/go-vue-selection/service"
)

func Setup(conf *config.AppConfig, db *gorm.DB) *gin.Engine {
	if conf.Mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode) // gin设置成发布模式
	}
	r := gin.New()

	// 装配各层依赖（依赖注入，去包级单例）
	attrRepo := repository.NewAttrRepo(db)
	categoryRepo := repository.NewCategoryRepo(db)
	menuRepo := repository.NewMenuRepo(db)
	roleRepo := repository.NewRoleRepo(db)
	skuRepo := repository.NewSkuRepo(db)
	spuRepo := repository.NewSpuRepo(db)
	trademarkRepo := repository.NewTrademarkRepo(db)
	userRepo := repository.NewUserRepo(db)

	// JWT 实例：密钥与有效期从配置注入，避免包内直接依赖全局 viper
	jwtSvc := jwt.NewJWT(viper.GetString("auth.jwt_secret"), viper.GetInt("auth.jwt_expire"))

	userHandler := handler.NewUserHandler(service.NewUserService(userRepo, roleRepo, jwtSvc))
	roleHandler := handler.NewRoleHandler(service.NewRoleService(roleRepo))
	menuHandler := handler.NewMenuHandler(service.NewMenuService(menuRepo))
	trademarkHandler := handler.NewTrademarkHandler(service.NewTrademarkService(trademarkRepo))
	categoryHandler := handler.NewCategoryHandler(service.NewCategoryService(categoryRepo))
	attrHandler := handler.NewAttrHandler(service.NewAttrService(attrRepo))
	spuHandler := handler.NewSpuHandler(service.NewSpuService(spuRepo))
	skuHandler := handler.NewSkuHandler(service.NewSkuService(skuRepo, spuRepo))
	fileHandler := handler.NewFileHandler(conf.Static.Path)

	// 图片路由设置
	r.MaxMultipartMemory = 4 << 20 // 4 MiB
	r.Static("/static", conf.Static.Path)

	//r.Use(middlewares.GinLogger(), logger.GinRecovery(true), middlewares.RateLimitMiddleware(2*time.Second, 1))
	r.Use(middlewares.GinLogger(), gin.Recovery(), middlewares.Cors())

	// swag route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.DocExpansion("none")))

	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "I'm OK!")
	})
	// 权限管理路由（登录接口无需 JWT，其余需 JWT 认证）
	adminAclGroup := r.Group("/admin/acl")
	userHandler.RegisterPublicRoutes(adminAclGroup) // 登录在 JWT 中间件之前注册
	adminAclGroup.Use(middlewares.JWTAuthMiddleware(jwtSvc))

	roleHandler.RegisterRoutes(adminAclGroup)
	menuHandler.RegisterRoutes(adminAclGroup)
	userHandler.RegisterRoutes(adminAclGroup)

	// 商品管理路由（整体需要 JWT 认证）
	adminProductGroup := r.Group("/admin/product")
	adminProductGroup.Use(middlewares.JWTAuthMiddleware(jwtSvc))

	categoryHandler.RegisterRoutes(adminProductGroup)
	trademarkHandler.RegisterRoutes(adminProductGroup)
	attrHandler.RegisterRoutes(adminProductGroup)
	spuHandler.RegisterRoutes(adminProductGroup)
	skuHandler.RegisterRoutes(adminProductGroup)
	fileHandler.RegisterRoutes(adminProductGroup)

	r.NoRoute(func(c *gin.Context) {
		result.Error(c, result.CodeNoRoute)
	})
	return r
}
