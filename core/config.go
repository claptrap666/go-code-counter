package core

var CurrentConfig *Config = &Config{}

type Config struct {
	Server *ServerConfig `mapstructure:"server"`
	DB     *DBConfig     `mapstructure:"db"`
	Oauth  *OauthConfig  `mapstructure:"oauth"`
}

type ServerConfig struct {
	Port   int    `mapstructure:"port"`
	Secret string `mapstructure:"secret"`
}

type DBConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
}

type OauthConfig struct {
	Provider     string `mapstructure:"provider"`
	ClientID     string `mapstructure:"client-id"`
	ClientSecret string `mapstructure:"client-secret"`
	RedirectURL  string `mapstructure:"redirect-url"`
}
