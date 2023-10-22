package main

import (
	. "github.com/Rayson9418/cdc/logger"
	"go.uber.org/zap"

	"examples/mongo"
	"examples/mysql"
	"examples/options"
)

func main() {
	if err := options.Init(); err != nil {
		Logger.Fatal("init cdc options err", zap.Error(err))
	}

	errChn := make(chan error, 2)
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

	for err := range errChn {
		Logger.Error("demo monitor err", zap.Error(err))
	}

	// TODO: translate chinese comment into english
	// TODO: we should run only one app to watch all database, mysql is ok, but watching mongo should be updated to client-level.
	// TODO: syncer needs to add example case and test.
	// TODO: next step, put data from monitoring into kafka
	// TODO: next step, add benthos
}
