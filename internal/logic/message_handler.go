package logic

import (
	"context"
	"fmt"

	"github.com/eatmoreapple/openwechat"
	"github.com/sashabaranov/go-openai"
)

func (t *TalkBotImpl) Handle(ctx context.Context, msg *openwechat.Message) error {
	if msg.IsText() {
		log.Info(fmt.Sprintf("text message:%s", msg.Content))
		resp, err := t.OpenaiSvr.ChatCompletion(ctx, generateChatMessages(msg.Content))
		if err != nil {
			msg.ReplyText("something bad happened please talk me later")
			return err
		}
		msg.ReplyText(resp.Choices[0].Message.Content)
	}
	return nil
}

func generateChatMessages(prompt string) []openai.ChatCompletionMessage {
	return []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleUser,
			Content: prompt,
		},
	}
}
