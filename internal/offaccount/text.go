package offaccount

import (
	"context"
	"talk_bot/internal/log"

	"github.com/silenceper/wechat/v2/officialaccount/message"
)

func (t *TalkBotImpl) onTextMessage(ctx context.Context, msg *message.MixMessage) *message.Reply {
	log.Infof("text message:%s", msg.Content)
	resp, err := t.OpenaiSvr.ChatCompletion(ctx, generateChatMessages(msg.GetOpenID(), msg.Content, false))
	if err != nil {
		return newTextReply("something bad happened please talk me later")
	}
	reply := resp.Choices[0].Message.Content
	saveMessageContext(msg.GetOpenID(), msg.Content, reply)
	return newTextReply(reply)
}

func newTextReply(reply string) *message.Reply {
	return &message.Reply{
		MsgType: message.MsgTypeText,
		MsgData: message.NewText(reply),
	}
}
