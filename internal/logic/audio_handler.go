package logic

import (
	"context"
	"fmt"
	"talk_bot/internal/service/openai"

	"github.com/eatmoreapple/openwechat"
)

func (t *TalkBotImpl) onAudioMessage(ctx context.Context, msg *openwechat.Message) error {
	audioResp, err := msg.GetVoice()
	if err != nil {
		_, _ = msg.ReplyText("I apologize, but i can not reach your audio,please talk to me later!")
		return fmt.Errorf("get voice failed:%v", err)
	}
	defer audioResp.Body.Close()

	// transcript user audio to text
	text, err := t.OpenaiSvr.SpeechToText(ctx, &openai.AudioRequest{
		Reader: audioResp.Body,
	})
	if err != nil {
		_, _ = msg.ReplyText("I apologize, but there is something went wrong,please talk to me later!")
		return fmt.Errorf("transcripte audio failed:%v", err)
	}

	// get the answer of user's audit
	resp, err := t.OpenaiSvr.ChatCompletion(ctx, generateChatMessages(msg.FromUserName, text))
	reply := resp.Choices[0].Message.Content
	_, _ = msg.ReplyText(reply)
	saveMessageContext(msg.FromUserName, text, reply)
	return nil
}
