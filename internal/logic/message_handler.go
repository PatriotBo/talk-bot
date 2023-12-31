package logic

import (
	"context"
	"fmt"
	"talk_bot/internal/log"
	"time"

	"github.com/eatmoreapple/openwechat"
	"github.com/patrickmn/go-cache"
	"github.com/sashabaranov/go-openai"
)

const (
	roundLimit      = 10
	roundExpireTime = 30 * time.Minute
)

// history rounds of user conversation with bot
var conversationRoundHistory *cache.Cache

func init() {
	conversationRoundHistory = cache.New(roundExpireTime, roundExpireTime)
}

func (t *TalkBotImpl) Handle(ctx context.Context, msg *openwechat.Message) error {
	switch msg.MsgType {
	case openwechat.MsgTypeText:
		return t.onTextMessage(ctx, msg)
	case openwechat.MsgTypeVoice:
		return t.onAudioMessage(ctx, msg)
	default:
		return fmt.Errorf("not supported message type:%v", msg.MsgType)
	}
}

func (t *TalkBotImpl) onTextMessage(ctx context.Context, msg *openwechat.Message) error {
	log.Infof("text message:%s", msg.Content)
	resp, err := t.OpenaiSvr.ChatCompletion(ctx, generateChatMessages(msg.FromUserName, msg.Content))
	if err != nil {
		_, _ = msg.ReplyText("something bad happened please talk me later")
		return err
	}
	reply := resp.Choices[0].Message.Content
	_, _ = msg.ReplyText(reply)
	saveMessageContext(msg.FromUserName, msg.Content, reply)
	return nil
}

func saveMessageContext(userID, prompt, reply string) {
	currentRound := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleUser,
			Content: prompt,
		},
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: reply,
		},
	}
	hisRoundsI, ok := conversationRoundHistory.Get(userID)
	if !ok {
		conversationRoundHistory.Set(userID, currentRound, cache.DefaultExpiration)
		return
	}

	rounds := hisRoundsI.([]openai.ChatCompletionMessage)
	rounds = append(rounds, currentRound...)
	// only save the latest N rounds of conversation
	if len(rounds) > roundLimit {
		rounds = rounds[len(rounds)-roundLimit:]
	}
	conversationRoundHistory.Set(userID, rounds, cache.DefaultExpiration)
}

func generateChatMessages(userID, prompt string) []openai.ChatCompletionMessage {
	promptMessage := openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: prompt,
	}

	hisRoundsI, ok := conversationRoundHistory.Get(userID)
	if !ok {
		return []openai.ChatCompletionMessage{promptMessage}
	}

	rounds := hisRoundsI.([]openai.ChatCompletionMessage)
	return append(rounds, promptMessage)
}
