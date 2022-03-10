package router

import (
	"github.com/gin-gonic/gin"
	"shop-api/user-web/api"
)

func InitUserRouter(Router *gin.RouterGroup) {
	UserRouter := Router.Group("user")

	{
		UserRouter.GET("list", api.GetUserList)
		UserRouter.POST("pwd_login", api.PassWordLogin)

	}
}
