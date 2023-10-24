package syncer

import (
	"time"

	cdcstore "github.com/Rayson9418/cdc/store"
	"github.com/Rayson9418/cdc/syncer"
	"go.uber.org/zap"

	. "examples/logger"
	"examples/mongo"
	"examples/mysql/models"
)

type Demo1Syncer struct {
	syncer.DummySyncer
}

func (s *Demo1Syncer) SetPosStore(store cdcstore.SyncerPosInterface) {
	s.SyncerPosInterface = store
}

func (s *Demo1Syncer) QueryCount(start, end time.Time) (int64, error) {
	return models.QueryDemo1CountByTime(start, end)
}

func (s *Demo1Syncer) QueryData(offset, limit int, start, end time.Time) (interface{}, int64, error) {
	return models.QueryDemo1DataByTime(offset, limit, start, end)
}

func (s *Demo1Syncer) Sink(data interface{}) error {
	// Do your business logic
	demo1ModelList := data.([]*models.Demo1Model)

	for _, d := range demo1ModelList {
		Logger.Info("============= sink data =============", zap.Any("data", d))
	}
	return nil
}

func (s *Demo1Syncer) InitData() error {
	return nil
}

func (s *Demo1Syncer) InitialPos() (int64, int64, int64) {
	return time.Now().AddDate(-1, 0, 0).Unix(), time.Now().Unix(), 0
}

func (s *Demo1Syncer) NextPos(start, end, pos int64) (int64, int64, int64) {
	start = end
	end = end + 30
	return start, end, 0
}

func (s *Demo1Syncer) Interval() error {
	time.Sleep(time.Duration(30) * time.Second)
	return nil
}

func NewDemo1Syncer() *Demo1Syncer {
	// New position store
	coll := mongo.GetClient().Database("record").Collection("position")
	store := cdcstore.NewSyncerPosMgoStore(coll, "syncer:demo:demo1", mongo.GetTimeoutLimit())

	// New syncer
	s := new(Demo1Syncer)
	s.SyncerName = "demo.demo1"
	s.SetPosStore(store)

	return s
}
