// Package handler : pm private message
package handler

import (
	"github.com/eatmoreapple/openwechat"
	"log"
	"strings"
	"wechat_llm/llm/openai"
)

var _ MessageHandlerInterface = (*PrivateMessageHandler)(nil)

// PrivateMessageHandler UserMessageHandler 私聊消息处理
type PrivateMessageHandler struct {
}

// handle 处理消息
func (g *PrivateMessageHandler) handle(msg *openwechat.Message) error {
	if msg.IsText() {
		return g.ReplyText(msg)
	}
	return nil
}

// NewPrivateMessageHandler NewPrivateMessageHandler 创建私聊处理器
func NewPrivateMessageHandler() MessageHandlerInterface {
	return &PrivateMessageHandler{}
}

// ReplyText 发送文本消息到群
func (g *PrivateMessageHandler) ReplyText(msg *openwechat.Message) error {
	// 接收私聊消息
	sender, err := msg.Sender()
	log.Printf("Received User %v Text Msg : %v", sender.NickName, msg.Content)

	// 向GPT发起请求
	requestText := strings.TrimSpace(msg.Content)
	requestText = strings.Trim(msg.Content, "\n")
	reply, err := openai.Completions(requestText)
	if err != nil {
		log.Printf("gtp request error: %v \n", err)
		msg.ReplyText("As a AI BOT, I am still learning. Please wait for a while.")
		return err
	}
	if reply == "" {
		return nil
	}

	// 回复用户
	reply = strings.TrimSpace(reply)
	reply = strings.Trim(reply, "\n")
	_, err = msg.ReplyText(reply)
	if err != nil {
		log.Printf("response user error: %v \n", err)
	}
	return err
}
