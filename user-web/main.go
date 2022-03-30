package main

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"shop-api/user-web/global"
	"shop-api/user-web/initialize"
	"shop-api/user-web/utils"
	myvalidator "shop-api/user-web/validator"
)

func main() {

	//1. 初始化Logger
	initialize.InitLogger()
	//2. 初始化配置文件
	initialize.InitConfig()
	//3. 初始化Router
	Router := initialize.Routers()
	//4. 初始化翻译
	if err := initialize.InitTrans("zh"); err != nil {
		panic(err)
	}
	//5. 初始化srv连接
	initialize.InitSrvConn()

	viper.AutomaticEnv()

	// 本地开发环境，固定端口，线上自动获取端口
	debug := viper.GetBool("GOPATH")
	if !debug {
		port, err := utils.GetFreePort()
		if err == nil {
			global.ServerConfig.Port = port
		}
	}
	// 注册验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("mobile", myvalidator.ValidateMobile)
		err := v.RegisterTranslation("mobile", global.Trans, func(ut ut.Translator) error {
			return ut.Add("mobile", "{0} 非法手机号码!", true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("mobile", fe.Field())
			return t
		})
		if err != nil {
			return
		}
	}

	zap.S().Debugf("启动服务器，端口：  %d", global.ServerConfig.Port)
	err := Router.Run(fmt.Sprintf(":%d", global.ServerConfig.Port))
	if err != nil {
		zap.S().Panic("启动失败", err.Error())
	}
}
