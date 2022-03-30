package config

type UserSrvConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
	Name string `mapstructure:"name" json:"name"`
}

type ServerConfig struct {
	Name          string        `mapstructure:"name" json:"name"`
	Port          int           `mapstructure:"port" json:"port"`
	UserSrcConfig UserSrvConfig `mapstructure:"user_srv" json:"user_srv"`
	JWTInfo       JWTConfig     `mapstructure:"jwt" json:"jwt"`
	AliSmsInfo    AliSmsConfig  `mapstructure:"sms" json:"sms"`
	RedisInfo     RedisConfig   `mapstructure:"redis" json:"redis"`
	ConsulInfo    ConsulConfig  `mapstructure:"consul" json:"consul"`
}

type JWTConfig struct {
	SigningKey string `mapstructure:"key" json:"key"`
}

type AliSmsConfig struct {
	ApiKey    string `mapstructure:"key" json:"key"`
	ApiSecret string `mapstructure:"secret" json:"secret"`
}

type RedisConfig struct {
	Host   string `mapstructure:"host" json:"host"`
	Port   int    `mapstructure:"port" json:"port"`
	Expire int    `mapstructure:"expire" json:"expire"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type NacosConfig struct {
	Host      string `mapstructure:"host"`
	Port      uint64 `mapstructure:"port"`
	Namespace string `mapstructure:"namespace"`
	User      string `mapstructure:"user"`
	Password  string `mapstructure:"password"`
	DataId    string `mapstructure:"dataid"`
	Group     string `mapstructure:"group"`
}
