package global

import (
	ut "github.com/go-playground/universal-translator"
	"shop-api/user-web/config"
	"shop-api/user-web/proto"
)

var (
	ServerConfig *config.ServerConfig = &config.ServerConfig{}
	Trans        ut.Translator

	UserSrvClient proto.UserClient
)
