package offaccount

import (
	"context"
	"fmt"
	"os"
	"talk_bot/internal/log"
	msgclient "talk_bot/internal/message"
	"talk_bot/internal/service/googlecloud"
	"talk_bot/internal/service/openai"

	"github.com/silenceper/wechat/v2/officialaccount/message"
)

const errReply = "I apologize, but there is something went wrong,please talk to me later!"

func (t *TalkBotImpl) onVoiceMessage(ctx context.Context, msg *message.MixMessage) *message.Reply {
	filepath, err := t.MsgClient.GetVoice(ctx, msg)
	if err != nil {
		log.Errorf("get voice failed:%v", err)
		return newTextReply("I apologize, but i can not reach your voice,please talk to me later!")
	}

	log.Infof("onVoiceMessage filepath:%s", filepath)

	// transcript user audio to text
	text, err := t.OpenaiSvr.SpeechToText(ctx, &openai.AudioRequest{
		FilePath: filepath,
	})
	if err != nil {
		log.Errorf("speech to text failed:%v", err)
		return newTextReply(errReply)
	}

	// get the answer of user's audit
	resp, err := t.OpenaiSvr.ChatCompletion(ctx, generateChatMessages(msg.GetOpenID(), text, true))
	reply := resp.Choices[0].Message.Content
	saveMessageContext(msg.GetOpenID(), text, reply)

	log.Infof("onVoiceMessage reply:%s", reply)

	// transcript text answer to speech
	filename, err := t.textToSpeech(ctx, msg.GetOpenID(), reply)
	if err != nil {
		log.Errorf("text to speech failed:%v", err)
		return newTextReply(errReply)
	}

	vr, err := t.voiceReply(ctx, filename)
	if err != nil {
		log.Errorf("voice reply failed:%v", err)
		return newTextReply(errReply)
	}
	return vr
}

func (t *TalkBotImpl) voiceReply(ctx context.Context, filename string) (*message.Reply, error) {
	mediaID, err := t.MsgClient.UploadVoice(ctx, filename)
	if err != nil {
		return nil, err
	}

	return newVoiceReply(mediaID), nil
}

func (t *TalkBotImpl) textToSpeech(ctx context.Context, username string, text string) (string, error) {
	by, err := t.TTSSvr.TextToSpeech(ctx, googlecloud.TTSRequest{Content: text, Language: googlecloud.LanguageEnUs})
	if err != nil {
		return "", err
	}
	path := msgclient.VoicePath(username, "reply", msgclient.MP3)
	if err = os.WriteFile(path, by, 0644); err != nil {
		return "", fmt.Errorf("write voice file :%v", err)
	}
	return path, nil
}

func newVoiceReply(mediaID string) *message.Reply {
	return &message.Reply{
		MsgType: message.MsgTypeVoice,
		MsgData: message.NewVoice(mediaID),
	}
}
