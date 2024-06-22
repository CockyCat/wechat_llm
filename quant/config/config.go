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

type OkxApiConf struct {
	ApiHost     string
	ApiKey      string
	SecretKey   string
	Passphrase  string
	IsSimulated bool
}

var config *OkxApiConf
var once sync.Once

// LoadOKXConfig 加载配置
func LoadOKXConfig() *OkxApiConf {
	once.Do(func() {
		config = &OkxApiConf{}

		projectRoot, err := getProjectRoot()

		fmt.Println("projectRoot:", projectRoot)
		if err != nil {
			log.Fatalf("Error reading project root: %v", err)
		}

		viper.SetConfigName("config")
		viper.SetConfigType("yaml")

		viper.AddConfigPath(filepath.Join(projectRoot, "quant/config"))

		fmt.Println(filepath.Join(projectRoot, "quant/config"))

		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("Error reading config file: %v", err)
		}
		config.ApiHost = viper.GetString("OkxApiConf.ApiHost")
		config.ApiKey = viper.GetString("OkxApiConf.ApiKey")
		config.SecretKey = viper.GetString("OkxApiConf.SecretKey")
		config.Passphrase = viper.GetString("OkxApiConf.Passphrase")
		config.IsSimulated = viper.GetBool("OkxApiConf.IsSimulated")
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
