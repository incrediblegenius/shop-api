package router

import (
	"github.com/gin-gonic/gin"
	"shop-api/user-web/api"
	"shop-api/user-web/middlewares"
)

func InitUserRouter(Router *gin.RouterGroup) {
	UserRouter := Router.Group("user")

	{
		UserRouter.GET("list", middlewares.JWTAuth(), middlewares.IsAdminAuth(), api.GetUserList)
		UserRouter.POST("pwd_login", api.PassWordLogin)
		//UserRouter.GET("list", api.GetUserList)
		UserRouter.POST("register", api.Register)
	}
}
