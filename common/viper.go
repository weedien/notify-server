package common

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type GlobalConfig struct {
	// MySQL数据库配置
	MySQL struct {
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
	}

	// 邮件服务配置
	EMail struct {
		Sender     string
		Password   string
		SMTPServer string
	}
}

var CONFIG GlobalConfig

func InitViper() *viper.Viper {
	v := viper.New()
	v.SetConfigFile("config.toml")
	v.SetConfigType("toml")
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	v.WatchConfig()

	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
		// 验证配置文件合法性
		if err = v.Unmarshal(&CONFIG); err != nil {
			fmt.Println(err)
		}
	})
	if err = v.Unmarshal(&CONFIG); err != nil {
		fmt.Println(err)
	}
	return v
}
