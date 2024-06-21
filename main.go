package main

import (
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"log"
	"wechat_llm/handler"
)

func main() {
	bot := openwechat.DefaultBot(openwechat.Desktop) // 桌面模式

	// 注册消息处理函数
	//bot.MessageHandler = func(msg *openwechat.Message) {
	//	//再此监听消息并回复消息
	//
	//	if msg.IsText() && msg.Content == "ping" {
	//		text, err := msg.ReplyText("pong")
	//		if err != nil {
	//			return
	//		}
	//		fmt.Println(text)
	//	}
	//}
	bot.MessageHandler = handler.Handler

	// 注册登陆二维码回调
	bot.UUIDCallback = openwechat.PrintlnQrcodeUrl

	// 创建热存储容器对象
	reloadStorage := openwechat.NewFileHotReloadStorage("storage.json")
	// 执行热登录
	err := bot.HotLogin(reloadStorage)
	if err != nil {
		if err = bot.Login(); err != nil {
			log.Printf("login error: %v \n", err)
			return
		}
	}

	// 登陆
	//if err := bot.Login(); err != nil {
	//	fmt.Println(err)
	//	return
	//}

	// 获取登陆的用户
	self, err := bot.GetCurrentUser()
	if err != nil {
		fmt.Println(err)
		return
	}

	// 获取所有的好友
	friends, err := self.Friends()
	fmt.Println("friends: ", friends, err)
	for _, friend := range friends {
		fmt.Println("NickName: ", friend.NickName)
		fmt.Println("AppAccountFlag: ", friend.AppAccountFlag)
		fmt.Println("City: ", friend.City)
		fmt.Println("RemarkName: ", friend.RemarkName)

		fmt.Println("===================")
	}

	// 获取所有的群组
	groups, err := self.Groups()
	fmt.Println("groups:", groups, err)

	// 阻塞主goroutine, 直到发生异常或者用户主动退出
	err = bot.Block()
	if err != nil {
		return
	}
}
