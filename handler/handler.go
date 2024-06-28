package handler

import (
	"log"
	"wechat_llm/config"

	"github.com/eatmoreapple/openwechat"
)

// MessageHandlerInterface message handler interface
type MessageHandlerInterface interface {
	handle(*openwechat.Message) error
	ReplyText(*openwechat.Message) error
}

type Type string

const (
	GroupHandler   = "group"
	PrivateHandler = "private"
)

// handlers 所有消息类型类型的处理器
var handlers map[Type]MessageHandlerInterface

func init() {
	handlers = make(map[Type]MessageHandlerInterface)
	handlers[GroupHandler] = NewGroupMessageHandler()
	handlers[PrivateHandler] = NewPrivateMessageHandler()
}

func Handler(msg *openwechat.Message) {
	log.Printf("hadler Received msg : %v", msg.Content)
	// group message
	if msg.IsSendByGroup() {
		err := handlers[GroupHandler].handle(msg)
		if err != nil {
			return
		}
		return
	}

	// request friend
	if msg.IsFriendAdd() {
		if config.LoadConfig().AutoPass {
			_, err := msg.Agree("HI Guy, IM an AI BOT，You can ask me anything, Are you ready, let got it!")
			if err != nil {
				log.Fatalf("add friend agree error : %v", err)
				return
			}
		}
	}

	// pm
	err := handlers[PrivateHandler].handle(msg)
	if err != nil {
		return
	}

	if msg.IsJoinGroup() {
		handler := GroupMessageHandler{}
		handler.Welcome(msg)
	}

}
