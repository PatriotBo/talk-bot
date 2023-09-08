package logic

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"talk_bot/internal/service/googlecloud"
	"talk_bot/internal/service/openai"
	"time"

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

	_, err = t.textToSpeech(ctx, msg.FromUserName, reply)
	if err != nil {
		_, _ = msg.ReplyText("I apologize, but there is something went wrong,please talk to me later!")
		return fmt.Errorf("text to speech failed:%v", err)
	}

	return nil
}

func (t *TalkBotImpl) textToSpeech(ctx context.Context, username string, text string) (io.Reader, error) {
	by, err := t.TTSSvr.TextToSpeech(ctx, googlecloud.TTSRequest{Content: text, Language: googlecloud.LanguageEnUs})
	if err != nil {
		fmt.Printf("textToSpeech failed err:%v \n", err)
		return nil, err
	}
	filename := fmt.Sprintf("%s_%d.mp3", username, time.Now().Unix())
	if err = os.WriteFile(filename, by, 0644); err != nil {
		fmt.Printf("textToSpeech write file err:%v \n", err)
		return nil, err
	}
	bytes.NewReader(by)
	return bytes.NewReader(by), nil
}
