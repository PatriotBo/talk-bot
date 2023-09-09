package offaccount

import (
	"net/http"

	"talk_bot/internal/conf"
	msgclient "talk_bot/internal/message"
	"talk_bot/internal/service/googlecloud"
	"talk_bot/internal/service/openai"

	"github.com/gin-gonic/gin"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/officialaccount"
	offconfig "github.com/silenceper/wechat/v2/officialaccount/config"
)

type TalkBotImpl struct {
	oa        *officialaccount.OfficialAccount
	OpenaiSvr openai.Service
	TTSSvr    googlecloud.Service
	MsgClient *msgclient.Message
}

func NewTalkBotImpl() *TalkBotImpl {
	oa := newWechatOfficialAccount()
	return &TalkBotImpl{
		oa:        oa,
		OpenaiSvr: openai.New(conf.GetConfig().OpenAI),
		TTSSvr:    googlecloud.NewTTS(),
		MsgClient: msgclient.NewMessage(oa),
	}
}

func newWechatOfficialAccount() *officialaccount.OfficialAccount {
	config := &offconfig.Config{
		AppID:          conf.GetWechatConfig().AppID,
		AppSecret:      conf.GetWechatConfig().AppSecret,
		Token:          conf.GetWechatConfig().Token,
		EncodingAESKey: conf.GetWechatConfig().EncodingAESKey,
		Cache:          cache.NewMemory(), // 使用本地缓存 保存 token

	}

	return officialaccount.NewOfficialAccount(config)
}

// Run start to service
func (t *TalkBotImpl) Run() {
	e := gin.Default()

	e.POST("/wx", func(ctx *gin.Context) {
		t.Handle(ctx)
	})

	server := &http.Server{
		Addr:    ":80",
		Handler: e,
	}
	// 启动HTTPS服务器
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
