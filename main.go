package main

import (
	"fmt"
	"talk_bot/internal/logic"
	"talk_bot/internal/offaccount"

	"github.com/eatmoreapple/openwechat"
)

func main() {
	// 网页版微信聊天机器人
	//runWebWechat()

	// 公众号（订阅号）聊天机器人
	imp := offaccount.NewTalkBotImpl()
	imp.Run()
}

// 网页版微信不支持语音功能
func runWebWechat() {
	bot := openwechat.DefaultBot(openwechat.Desktop) // 桌面模式

	impl := logic.NewTalkBotImpl()
	// 注册消息处理函数
	bot.MessageHandler = func(msg *openwechat.Message) {
		if msg.IsSendByGroup() {
			return
		}
		if err := impl.Handle(msg.Context(), msg); err != nil {
			fmt.Printf("handle message failed:%v \n", err)
		}
	}

	// 注册登陆二维码回调
	bot.UUIDCallback = openwechat.PrintlnQrcodeUrl

	// 登陆
	if err := bot.Login(); err != nil {
		fmt.Println(err)
		return
	}

	// 阻塞主goroutine, 直到发生异常或者用户主动退出
	bot.Block()
}
