package main

import (
	cdclogger "github.com/Rayson9418/cdc/logger"
	"go.uber.org/zap"

	. "examples/logger"
	"examples/mongo"
	"examples/mysql"
	"examples/options"
	"examples/syncer"
)

func main() {
	cdclogger.SetLogger(Logger)

	if err := options.Init(); err != nil {
		Logger.Fatal("init cdc options err", zap.Error(err))
	}

	if err := mongo.InitClient(); err != nil {
		Logger.Fatal("init mongo client err", zap.Error(err))
	}

	if err := mysql.InitClient(); err != nil {
		Logger.Fatal("init mysql client err", zap.Error(err))
	}

	errChn := make(chan error, 5)
	go func() {
		if err := mysql.DemoMonitor(); err != nil {
			errChn <- err
		}
	}()

	go func() {
		if err := mongo.DemoMonitor(); err != nil {
			errChn <- err
		}
	}()

	go func() {
		if err := syncer.DemoSyncOnce(); err != nil {
			errChn <- err
		}
	}()

	go func() {
		if err := syncer.DemoSyncOnTime(); err != nil {
			errChn <- err
		}
	}()

	go func() {
		if err := syncer.DemoSyncAlways(); err != nil {
			errChn <- err
		}
	}()

	for err := range errChn {
		Logger.Error("demo monitor err", zap.Error(err))
	}
}
