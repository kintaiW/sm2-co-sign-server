package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Auth     AuthConfig     `mapstructure:"auth"`
	Log      LogConfig      `mapstructure:"log"`
	Admin    AdminConfig    `mapstructure:"admin"`
}

type AdminConfig struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type ServerConfig struct {
	Port         int           `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

type DatabaseConfig struct {
	Path string `mapstructure:"path"`
}

type AuthConfig struct {
	TokenExpire time.Duration `mapstructure:"token_expire"`
	MasterKey   string        `mapstructure:"master_key"`
}

type LogConfig struct {
	Level  string `mapstructure:"level"`
	Output string `mapstructure:"output"`
}

var AppConfig *Config

func Load(configPath string) error {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	AppConfig = &Config{}
	return viper.Unmarshal(AppConfig)
}
