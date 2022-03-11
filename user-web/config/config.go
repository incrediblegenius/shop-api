package config

type UserSrvConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type ServerConfig struct {
	Name          string        `mapstructure:"name"`
	Port          int           `mapstructure:"port"`
	UserSrcConfig UserSrvConfig `mapstructure:"user_srv"`
	JWTInfo       JWTConfig     `mapstructure:"jwt"`
}

type JWTConfig struct {
	SigningKey string `mapstructure:"key"`
}
