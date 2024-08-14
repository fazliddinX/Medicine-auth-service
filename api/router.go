package api

import (
	"auth-service/api/handler"
	"auth-service/api/middleware"
	"auth-service/pkg/config"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "auth-service/api/docs"
)

type Router interface {
	InitRouter()
	RunRouter(cf config.Config) error
}

func NewRouter(authHandler handler.AuthHandler) Router {
	router := gin.Default()
	return &routerImpl{router: router, handler: authHandler}
}

type routerImpl struct {
	router  *gin.Engine
	handler handler.AuthHandler
}

// @title Authenfication service
// @version 1.0
// @description server for siginIn or signUp
// @host localhost:8081
// @BasePath /auth
// @schemes http
func (r *routerImpl) InitRouter() {

	r.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	auth := r.router.Group("/auth")
	auth.POST("/register", r.handler.Register)
	auth.POST("/login", r.handler.Login)
	auth.POST("/admin", r.handler.AddAdmin)
	auth.POST("/refresh", middleware.GetAccessTokenMid(), r.handler.RefreshToken)
	auth.GET("/get-role", r.handler.GetRole)

}

func (r *routerImpl) RunRouter(cf config.Config) error {
	return r.router.Run(cf.HOST + cf.GIN_SERVER_PORT)
}
