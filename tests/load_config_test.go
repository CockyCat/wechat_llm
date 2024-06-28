package tests

import (
	"testing"
	"wechat_llm/config"
	"wechat_llm/llm/openai"
	"wechat_llm/macro"
)

func TestLoadConfig(t *testing.T) {
	apiKey := config.LoadConfig().ApiKey
	t.Log("apiKey is: ", apiKey)
}

func TestLoadConfig2(t *testing.T) {
	//apiKey := config2.LoadOKXConfig().ApiKey
	//t.Log("apiKey is: ", apiKey)

	marcoFRContent := macro.RunAndGetData("DGS10")
	t.Log(marcoFRContent)
}

func TestGPTChatProxy(t *testing.T) {
	reply, err := openai.GPTProxyChat("你好")
	if err != nil {
		t.Error(err)
	}
	t.Log(reply)

}

func TestRunAndGetData(t *testing.T) {
	marcoFRContent := macro.RunAndGetData("SOFR")
	t.Log(marcoFRContent)
}
