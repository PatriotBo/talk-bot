package logic

import (
	"crypto/tls"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

var log *zap.Logger

func init() {
	// 创建一个基本的生产配置的 zap 日志记录器
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	log = logger
}

type TalkBotImpl struct {
	Engine *gin.Engine
}

func NewTalkBotImpl() *TalkBotImpl {
	return &TalkBotImpl{}
}

// Run start to service
func (m *TalkBotImpl) Run() {
	e := gin.Default()

	// 设置HTTPS证书和密钥
	certFile := "../certs/cert.pem"
	keyFile := "../certs/key.pem"

	// 配置TLS
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	server := &http.Server{
		Addr:      ":443",
		Handler:   e,
		TLSConfig: tlsConfig,
	}
	// 启动HTTPS服务器
	if err := server.ListenAndServeTLS(certFile, keyFile); err != nil {
		panic(err)
	}
}
