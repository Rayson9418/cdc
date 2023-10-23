package mongo

import (
	cdcmongo "github.com/Rayson9418/cdc/mongo"
	"go.uber.org/zap"

	. "examples/logger"
)

type DemoHandler struct {
	cdcmongo.DummyStreamHandler
}

func (h *DemoHandler) OnChange(object *cdcmongo.StreamObject) error {
	Logger.Info("start to handle stream object", zap.Any("object", object.FullDocument))

	data := ParseDemo(object.FullDocument)
	Logger.Info("this is what you parse from stream object", zap.Any("data", data))

	return nil
}

func NewDemoHandler() *DemoHandler {
	handler := &DemoHandler{}
	handler.DbName = "demo"
	handler.CollName = "demo1"
	handler.OpTypes = []string{"insert", "update", "delete"}
	return handler
}
