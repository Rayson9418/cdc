package syncer

import (
	"errors"
	"fmt"
	"math"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"

	"github.com/Rayson9418/cdc/common"
	. "github.com/Rayson9418/cdc/logger"
	"github.com/Rayson9418/cdc/store"
)

const kSyncerPosFmt = "%d:%d:%d"

type DataSyncInterface interface {
	store.SyncerPosInterface
	// Name define the syncer name
	Name() string
	// SetPosStore set position store
	SetPosStore(store.SyncerPosInterface)
	// QueryCount function for total data count retrieval
	QueryCount(start, end time.Time) (int64, error)
	// QueryData function for data retrieval
	QueryData(offset, limit int, start, end time.Time) (interface{}, int64, error)
	// Sink function for data synchronization logic
	Sink(data interface{}) error
	// InitData function for initialization logic before data synchronization logic
	InitData() error
	// InitialPos for the position during the initial synchronization
	InitialPos() (int64, int64, int64)
	// NextPos for the position of the next data synchronization
	NextPos(start, end, pos int64) (int64, int64, int64)
	// Interval for the synchronization interval
	Interval() error
}

type DummySyncer struct {
	SyncerName string
	store.SyncerPosInterface
}

func (s *DummySyncer) Name() string {
	return s.SyncerName
}

func (s *DummySyncer) SetPosStore(store store.SyncerPosInterface) {
	s.SyncerPosInterface = store
}

func (s *DummySyncer) QueryCount(start, end time.Time) (int64, error) {
	return 0, nil
}

func (s *DummySyncer) QueryData(offset, limit int, start, end time.Time) (interface{}, int64, error) {
	return 0, 0, nil
}

func (s *DummySyncer) Sink(interface{}) error {
	return nil
}

func (s *DummySyncer) InitData() error {
	return nil
}

func (s *DummySyncer) InitialPos() (int64, int64, int64) {
	return 0, 0, 0
}

func (s *DummySyncer) NextPos(start, end, pos int64) (int64, int64, int64) {
	return start, end, pos
}

func (s *DummySyncer) Interval() error {
	return nil
}

func StartSyncer(syncers ...DataSyncInterface) error {
	errChn := make(chan error, len(syncers))
	for _, s := range syncers {
		go func(s DataSyncInterface) {
			if err := syncAlways(s); err != nil {
				errChn <- err
			}
		}(s)
	}

	for err := range errChn {
		return err
	}
	return nil
}

func StartSyncerOnTime(syncers ...DataSyncInterface) error {

	errChn := make(chan error, len(syncers))
	for _, s := range syncers {
		go func(s DataSyncInterface) {
			if err := syncOnTime(s); err != nil {
				errChn <- err
			}
		}(s)
	}

	for err := range errChn {
		return err
	}
	return nil
}

// SyncOnce to immediately execute the synchronization logic once
func SyncOnce(s DataSyncInterface) error {
	return syncOnce(s)
}

func syncAlways(s DataSyncInterface) error {
	for {
		if err := syncOnce(s); err != nil {
			return err
		}
		// Wait for the next synchronization cycle
		if err := s.Interval(); err != nil {
			return err
		}
	}
}

func syncOnTime(s DataSyncInterface) error {
	for {
		hour := time.Now().Hour()
		if !(hour >= opt.StartHour && hour <= opt.EndHour) {
			time.Sleep(time.Duration(opt.Interval) * time.Hour)
			continue
		}
		if err := syncOnce(s); err != nil {
			return err
		}
		// Wait for the next synchronization cycle
		if err := s.Interval(); err != nil {
			return err
		}
	}
}

func syncOnce(s DataSyncInterface) error {
	Logger.Info("start syncer", zap.String("name", s.Name()))

	syncerPos, err := s.Pos()
	if err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			Logger.Warn("query latest sync time err", zap.Error(err))
			return err
		}
		// Execute initialization function during the initial synchronization
		if err = s.InitData(); err != nil {
			Logger.Error("init data failed", zap.String("sync_name", s.Name()))
			return err
		}
		// Obtain the position of the first execution
		syncerStartTimestamp, syncEndTimestamp, pos := s.InitialPos()
		syncerPos = &store.SyncerPos{
			SyncStartTimestamp: syncerStartTimestamp,
			SyncEndTimestamp:   syncEndTimestamp,
			Pos:                pos,
		}
	}

	if common.IsToday(syncerPos.LastSyncTime) {
		Logger.Info("syncer has already completed today, skip!!!",
			zap.String("sync_name", s.Name()),
			zap.Int64("run_time", syncerPos.LastSyncTime))
		return nil
	}

	return startSyncer(s, syncerPos)
}

func startSyncer(s DataSyncInterface, syncerPos *store.SyncerPos) error {
	var (
		syncStartTimestamp = syncerPos.SyncStartTimestamp
		syncEndTimestamp   = syncerPos.SyncEndTimestamp
		pos                = syncerPos.Pos
	)

	endDate := time.Unix(syncEndTimestamp, 0)
	startDate := time.Unix(syncStartTimestamp, 0)

	// Query total count of data
	count, err := s.QueryCount(startDate, endDate)
	if err != nil {
		Logger.Error("query data count err",
			zap.String("sync_name", s.Name()),
			zap.Error(err))
		return err
	}
	Logger.Info("query data count",
		zap.String("sync_name", s.Name()),
		zap.Int64("count", count),
		zap.Time("start", startDate),
		zap.Time("end", endDate))

	// Current date has been fully synchronized
	if count != 0 && count <= pos {
		Logger.Info("last sync complete, start new sync",
			zap.Int64("pos", pos),
			zap.String("sync_name", s.Name()))

		// Set the next synchronization time range after synchronization is complete
		syncStartTimestamp, syncEndTimestamp, pos = s.NextPos(syncStartTimestamp, syncEndTimestamp, pos)

		count, err = s.QueryCount(startDate, endDate)
		if err != nil {
			Logger.Error("query data count err",
				zap.String("sync_name", s.Name()),
				zap.Error(err))
			return err
		}
	}

	// Paging query
	queryTimes := math.Ceil(float64(count-pos) / float64(opt.BatchLimit))
	offset := int(pos)
	for i := 0; i < int(queryTimes); i++ {
		dataList, nums, err := s.QueryData(offset, opt.BatchLimit, startDate, endDate)
		if err != nil {
			Logger.Error("query data err", zap.String("sync_name", s.Name()),
				zap.Error(err))
			return err
		}

		offset += opt.BatchLimit
		if err = s.Sink(dataList); err != nil {
			Logger.Error("sink data err",
				zap.String("sync_name", s.Name()),
				zap.Error(err))
			return err
		}

		pos += nums
		err = s.Save(fmt.Sprintf(kSyncerPosFmt, syncStartTimestamp, syncEndTimestamp, pos))
		if err != nil {
			return err
		}
	}

	Logger.Info("syncer end, wait next time", zap.String("name", s.Name()))
	return nil
}
