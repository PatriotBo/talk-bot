package log

import (
	"fmt"

	"go.uber.org/zap"
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

func Infof(format string, v any) {
	log.Info(fmt.Sprintf(format, v))
}

func Errorf(format string, v any) {
	log.Error(fmt.Sprintf(format, v))
}
func Debugf(format string, v any) {
	log.Debug(fmt.Sprintf(format, v))
}
