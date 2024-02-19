package util

import "go.uber.org/zap"

var globalLogger = MustNewLogger()

func GetGlobalLogger() *zap.Logger {
	return globalLogger
}
func MustNewLogger() *zap.Logger {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	return logger
}
