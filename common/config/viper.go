package config

import (
	"fmt"
	"github.com/sirupsen/logrus"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type GlobalConfig struct {
	// MySQL数据库配置
	Database struct {
		Host     string
		Port     int
		Username string
		Password string
		Database string
		Params   string
	}

	// http服务配置
	Server struct {
		Port int
		Env  string
	}

	// 邮件服务配置
	EMail struct {
		Sender     string
		Password   string
		SMTPServer string
	}

	// 安全配置
	Security struct {
		CorsAllowedOrigins string
	}

	Log struct {
		Level      string
		Stdout     bool
		OutputPath string
		Format     string
		Rotate     bool
	}
}

var conf GlobalConfig

func InitViper() *viper.Viper {
	v := viper.New()
	v.SetConfigFile("../config.toml")
	v.SetConfigType("toml")
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	v.WatchConfig()

	v.OnConfigChange(func(e fsnotify.Event) {
		logrus.Info("Config file changed:", e.Name)
		if err = v.Unmarshal(&conf); err != nil {
			fmt.Println(err)
		}
	})
	if err = v.Unmarshal(&conf); err != nil {
		fmt.Println(err)
	}
	return v
}

func Config() GlobalConfig {
	if conf == (GlobalConfig{}) {
		InitViper()
	}
	return conf
}
