package logic

import (
	"talk_bot/internal/conf"
	"talk_bot/internal/service/googlecloud"
	"talk_bot/internal/service/openai"
)

type TalkBotImpl struct {
	OpenaiSvr openai.Service
	TTSSvr    googlecloud.Service
}

func NewTalkBotImpl() *TalkBotImpl {
	return &TalkBotImpl{
		OpenaiSvr: openai.New(conf.GetConfig().OpenAI),
		TTSSvr:    googlecloud.NewTTS(),
	}
}
