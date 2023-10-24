package mysql

import (
	cdcmysql "github.com/Rayson9418/cdc/mysql"
	cdcredis "github.com/Rayson9418/cdc/redis"
	cdcstore "github.com/Rayson9418/cdc/store"
	"go.uber.org/zap"

	. "examples/logger"
	"examples/options"
)

func DemoMonitor() error {
	// InitClient position store
	if err := cdcredis.InitClient(options.CdcOpt.Redis); err != nil {
		Logger.Fatal("init redis client with opt err", zap.Error(err))
	}
	store := cdcstore.NewBinlogRedisStore("mysql:binlog:pos")

	// New handler for specific collection
	handler := NewDemoHandler()

	// New row event monitor
	m, err := cdcmysql.NewRowEventMonitor(options.CdcOpt.Mysql)
	if err != nil {
		return err
	}
	// Set position store for monitor
	m.SetStore(store)
	// Add handlers for monitor
	if err = m.AddHandler(handler); err != nil {
		return err
	}
	return m.StartMonitor()
}
