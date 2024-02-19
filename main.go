package main

import (
	"context"

	"github.com/kcwillis/jackhenry/weather/server"
	"github.com/kcwillis/jackhenry/weather/util"
	"go.uber.org/zap"
)

func main() {
	err := server.Serve(context.Background())
	if err != nil {
		util.GetGlobalLogger().Fatal("server fail", zap.Error(err))
	}
}
