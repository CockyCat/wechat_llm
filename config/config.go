package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	//viper
	"github.com/spf13/viper"
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
		ApiKey := os.Getenv("OPENAI_API_KEY")
		if ApiKey != "" {
			config.ApiKey = ApiKey
		} else {
			//读取config/config.yaml配置文件中的OPENAI_API_KEY
			fmt.Println("读取config/config.yaml配置文件中的OPENAI_API_KEY")
			// 获取项目根目录
			projectRoot, err := getProjectRoot()

			fmt.Println("projectRoot:", projectRoot)
			if err != nil {
				log.Fatalf("Error reading project root: %v", err)
			}

			// 设置配置文件名（不带扩展名）
			viper.SetConfigName("config")
			viper.SetConfigType("yaml")

			viper.AddConfigPath(filepath.Join(projectRoot, "config"))

			if err := viper.ReadInConfig(); err != nil {
				log.Fatalf("Error reading config file: %v", err)
			}
			config.ApiKey = viper.GetString("OPENAI_API_KEY")
		}
		config.AutoPass = true
	})
	return config
}

// getProjectRoot 获取项目根目录
func getProjectRoot() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(currentDir, "config")); os.IsNotExist(err) {
			parentDir := filepath.Dir(currentDir)
			if parentDir == currentDir {
				// 已经到根目录仍未找到
				return "", fmt.Errorf("project root not found")
			}
			currentDir = parentDir
		} else {
			return currentDir, nil
		}
	}
}
