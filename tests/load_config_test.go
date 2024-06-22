package tests

import (
	"testing"
	"wechat_llm/config"
)

func TestLoadConfig(t *testing.T) {
	apiKey := config.LoadConfig().ApiKey
	t.Log("apiKey is: ", apiKey)
}
