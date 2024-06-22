package tests

import (
	"testing"
	"wechat_llm/config"
	config2 "wechat_llm/quant/config"
)

func TestLoadConfig(t *testing.T) {
	apiKey := config.LoadConfig().ApiKey
	t.Log("apiKey is: ", apiKey)
}

func TestLoadConfig2(t *testing.T) {
	apiKey := config2.LoadOKXConfig().ApiKey
	t.Log("apiKey is: ", apiKey)
}
