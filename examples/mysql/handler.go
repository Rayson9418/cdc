package mysql

import (
	cdcmysql "github.com/Rayson9418/cdc/mysql"
	"go.uber.org/zap"

	. "examples/logger"
)

type DemoHandler struct {
	cdcmysql.DummyBinLogHandler
}

func (h *DemoHandler) OnRow(event *cdcmysql.RowEventData) error {
	Logger.Info("start to handle event data", zap.String("table", event.Table))

	data := ParseDemo(event)
	Logger.Info("this is what you parse from RowEventData", zap.Any("data", data))

	return nil
}

func NewDemoHandler() *DemoHandler {
	handler := &DemoHandler{}
	handler.DbName = "demo"
	handler.TableName = "demo1"
	handler.Actions = []string{"insert", "update", "delete"}
	return handler
}
