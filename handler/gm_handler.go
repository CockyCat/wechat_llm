package handler

import (
	"log"
	"strings"
	"wechat_llm/llm/openai"
	"wechat_llm/macro"

	"github.com/eatmoreapple/openwechat"
)

var _ MessageHandlerInterface = (*GroupMessageHandler)(nil)

// GroupMessageHandler 群消息处理
type GroupMessageHandler struct {
}

// handle 处理消息
func (g *GroupMessageHandler) handle(msg *openwechat.Message) error {
	if msg.IsText() {
		return g.ReplyText(msg)
	}
	return nil
}

// NewGroupMessageHandler 创建群消息处理器
func NewGroupMessageHandler() MessageHandlerInterface {
	return &GroupMessageHandler{}
}

// ReplyText 发送文本消息到群
func (g *GroupMessageHandler) ReplyText(msg *openwechat.Message) error {
	// 接收群消息
	sender, err := msg.Sender()
	group := openwechat.Group{User: sender}
	log.Printf("Received Group %v Text Msg : %v", group.NickName, msg.Content)

	if msg.Content == "FEDFUNDS" || msg.Content == "SOFR" || msg.Content == "DGS10" {
		marcoFRContent := macro.RunAndGetData(msg.Content)
		_, err := msg.ReplyText(marcoFRContent)
		if err != nil {
			log.Printf("response group error: %v \n", err)
		}
		log.Printf("FR marco data have sended to Group.")
	}

	// 不是@的不处理
	if !msg.IsAt() {
		return nil
	}

	// 替换掉@文本，然后向GPT发起请求
	replaceText := "@" + sender.NickName
	requestText := strings.TrimSpace(strings.ReplaceAll(msg.Content, replaceText, ""))
	//reply, err := openai.Completions(requestText)
	reply, err := openai.GPTProxyChat(requestText)

	if err != nil {
		log.Printf("gtp request error: %v \n", err)
		_, err := msg.ReplyText("...")
		if err != nil {
			return err
		}
		return err
	}
	if reply == "" {
		return nil
	}

	// 获取@我的用户
	groupSender, err := msg.SenderInGroup()
	if err != nil {
		log.Printf("get sender in group error :%v \n", err)
		return err
	}

	// 回复@我的用户
	reply = strings.TrimSpace(reply)
	reply = strings.Trim(reply, "\n")
	atText := "@" + groupSender.NickName
	replyText := atText + reply
	_, err = msg.ReplyText(replyText)
	if err != nil {
		log.Printf("response group error: %v \n", err)
	}
	return err
}

func (g *GroupMessageHandler) Welcome(msg *openwechat.Message) {
	content := msg.Content
	parts := strings.Split(content, "加入了群聊")
	if len(parts) != 2 {
		log.Printf("content split error: %v", content)
		return
	}
	nickName := parts[0]
	log.Printf("new member %v join group", nickName)

	reply := "欢迎新成员" + nickName + "加入群聊, 请报上三围、身高、体重。"

	atText := "@" + nickName
	replyText := atText + reply
	_, err := msg.ReplyText(replyText)
	if err != nil {
		log.Printf("welcome new member error: %v", err)
	}
}
