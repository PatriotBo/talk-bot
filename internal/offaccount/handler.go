package offaccount

import (
	"fmt"
	"time"

	"talk_bot/internal/log"

	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"github.com/sashabaranov/go-openai"
	"github.com/silenceper/wechat/v2/officialaccount/message"
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

func (t *TalkBotImpl) Handle(ctx *gin.Context) {
	server := t.oa.GetServer(ctx.Request, ctx.Writer)
	server.SkipValidate(false) // 跳过请求合法性检查
	server.SetMessageHandler(func(msg *message.MixMessage) *message.Reply {
		switch msg.MsgType {
		case message.MsgTypeText:
			return t.onTextMessage(ctx, msg)
		case message.MsgTypeVoice:
			return t.onVoiceMessage(ctx, msg)
		default:
			return &message.Reply{
				MsgType: message.MsgTypeText,
				MsgData: message.NewText(fmt.Sprintf("暂不支持的消息类型 ：%s", msg.MsgType)),
			}
		}
	})
	if err := server.Serve(); err != nil {
		log.Errorf("server.Serve failed err:%v", err.Error())
		return
	}
	// 发送回复的消息
	if err := server.Send(); err != nil {
		log.Errorf("server.Send failed err:%v", err.Error())
		return
	}
	log.Infof("HandleMessage success %s", server.Token)
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

var wordLimitPrompt = "Please answer all of my questions in no more than 50 words."

func generateChatMessages(userID, prompt string, limit bool) []openai.ChatCompletionMessage {
	if limit {
		prompt = prompt + wordLimitPrompt
	}
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
