package config

type UserSrvConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
	Name string `json:"name"`
}

type ServerConfig struct {
	Name          string        `json:"name"`
	Port          int           `json:"port"`
	UserSrcConfig UserSrvConfig `json:"user_srv"`
	JWTInfo       JWTConfig     `json:"jwt"`
	AliSmsInfo    AliSmsConfig  `json:"sms"`
	RedisInfo     RedisConfig   `json:"redis"`
	ConsulInfo    ConsulConfig  `json:"consul"`
}

type JWTConfig struct {
	SigningKey string `json:"key"`
}

type AliSmsConfig struct {
	ApiKey    string `json:"key"`
	ApiSecret string `json:"secret"`
}

type RedisConfig struct {
	Host   string `json:"host"`
	Port   int    `json:"port"`
	Expire string `json:"expire"`
}

type ConsulConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}
