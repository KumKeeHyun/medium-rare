package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/KumKeeHyun/medium-rare/trend-service/controller"
	"github.com/KumKeeHyun/medium-rare/trend-service/dao/sql"
	"github.com/KumKeeHyun/medium-rare/trend-service/util"
	"github.com/KumKeeHyun/medium-rare/trend-service/util/erouter"
)

func main() {
	logger, err := util.BuildZapLogger()
	if err != nil {
		panic(err)
	}

	db, err := util.BuildMysqlConnection()
	if err != nil {
		panic(err)
	}

	rrr := sql.NewSqlReadRecordRepository(db)
	ec := controller.NewEventController(rrr, logger)

	er := erouter.NewEventRouter("trend", logger)
	er.SetHandler("read-article", ec.ReadArticle)

	if err := er.StartRouter(); err != nil {
		panic(err)
	}

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	<-sigterm
	logger.Info("stop")
	er.Stop()
}
