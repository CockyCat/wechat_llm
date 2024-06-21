package config

import (
	"os"
	"sync"
)

// Configuration 项目配置
type Configuration struct {
	// gtp apikey
	ApiKey string `json:"api_key"`
	// 自动通过好友
	AutoPass bool `json:"auto_pass"`
}

var config *Configuration
var once sync.Once

// LoadConfig 加载配置
func LoadConfig() *Configuration {
	once.Do(func() {
		// 从文件中读取
		config = &Configuration{}
		// 如果环境变量有配置，读取环境变量
		ApiKey := os.Getenv("ApiKey")
		if ApiKey != "" {
			config.ApiKey = ApiKey
		}
		config.AutoPass = true
	})
	return config
}
