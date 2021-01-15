package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/KumKeeHyun/medium-rare/trend-service/util"
	"github.com/KumKeeHyun/medium-rare/trend-service/util/erouter"
	"go.uber.org/zap"
)

func main() {
	logger, err := util.BuildZapLogger()
	if err != nil {
		panic(err)
	}

	cg, err := util.BuildConsumerGroup(logger)
	if err != nil {
		panic(err)
	}

	er := erouter.NewEventRouter(cg, logger)
	er.SetHandler("read-article", func(key, value []byte) {
		logger.Info("handle event",
			zap.String("topic", "read-article"),
			zap.ByteString("value", value))
	})
	er.SetHandler("create-user", func(key, value []byte) {
		logger.Info("handle event",
			zap.String("topic", "create-user"),
			zap.ByteString("value", value))
	})

	er.Start()

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	<-sigterm
	logger.Info("stop")
	er.Stop()
}
