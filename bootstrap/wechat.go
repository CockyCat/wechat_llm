package bootstrap

import (
	"fmt"
	"log"
	"wechat_llm/handler"

	"github.com/eatmoreapple/openwechat"
)

type Wechat struct{}

func (w Wechat) Run() {
	bot := openwechat.DefaultBot(openwechat.Desktop) // 桌面模式

	// 注册消息处理函数
	bot.MessageHandler = handler.Handler
	// 注册登陆二维码回调
	bot.UUIDCallback = openwechat.PrintlnQrcodeUrl

	// 登陆
	//if err := bot.Login(); err != nil {
	//	fmt.Println(err)
	//	return
	//}

	// 创建热存储容器对象	这种方式如果换了电脑的话，则捕获不到消息，换电脑运行后，记得删除下storage.json这个文件
	reloadStorage := openwechat.NewFileHotReloadStorage("storage.json")
	// 执行热登录
	err := bot.HotLogin(reloadStorage)
	if err != nil {
		if err = bot.Login(); err != nil {
			log.Printf("login error: %v \n", err)
			return
		}
	}

	// 获取登陆的用户
	self, err := bot.GetCurrentUser()
	if err != nil {
		fmt.Println(err)
		return
	}

	// 获取所有的好友
	//friends, err := self.Friends()
	//for _, friend := range friends {
	//	msg, err := friend.SendText("检查下是否还是好友")
	//	fmt.Println(fmt.Sprintf("检测好友 %v", friend.Alias))
	//
	//	fmt.Println()
	//
	//	isFriend := friend.IsFriend()
	//	if isFriend == false {
	//		fmt.Println(fmt.Sprintf("%v 不是好友，可以删除。", friend.Alias))
	//	}
	//
	//	if err != nil {
	//		fmt.Println("检测好友消息err:", err)
	//	}
	//}

	//fmt.Println(friends, err)

	// 获取所有的群组
	groups, err := self.Groups()
	fmt.Println(groups, err)

	// 阻塞主goroutine, 直到发生异常或者用户主动退出
	err = bot.Block()
	if err != nil {
		return
	}
}
