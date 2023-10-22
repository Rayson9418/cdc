package mongo

import (
	. "github.com/Rayson9418/cdc/logger"
	cdcmongo "github.com/Rayson9418/cdc/mongo"
	cdcredis "github.com/Rayson9418/cdc/redis"
	cdcstore "github.com/Rayson9418/cdc/store"
	"go.uber.org/zap"

	"examples/options"
)

func DemoMonitor() error {
	// Init position store
	if err := cdcredis.InitClient(options.CdcOpt.Redis); err != nil {
		Logger.Fatal("init redis client with opt err", zap.Error(err))
	}
	store := cdcstore.NewStreamRedisStore("mongo:oplog:pos")

	// New handler for specific collection
	handler := NewDemoHandler()

	// New mongo monitor
	m, err := cdcmongo.NewDefaultMonitor(options.CdcOpt.Mongo)
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
