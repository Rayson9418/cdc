package store

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"

	. "github.com/Rayson9418/cdc/logger"
)

type MongoPos struct {
	Value   string
	Version int64
}

type SyncerPosMgoStore struct {
	coll    *mongo.Collection
	PosKey  string
	Timeout time.Duration
}

func NewSyncerPosMgoStore(coll *mongo.Collection, key string) *SyncerPosMgoStore {
	return &SyncerPosMgoStore{
		PosKey: key,
		coll:   coll,
	}
}

func (s *SyncerPosMgoStore) getPosStr() (value string, version int64, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.Timeout)
	defer cancel()

	filter := bson.M{
		"key": s.PosKey,
	}

	pos := &MongoPos{}
	if err := s.coll.FindOne(ctx, filter).Decode(pos); err != nil {
		Logger.Warn("get pos str err", zap.String("key", s.PosKey),
			zap.Error(err))
		return "", 0, err
	}
	return pos.Value, pos.Version, nil
}

func (s *SyncerPosMgoStore) setPosStr(value interface{}, unix ...int64) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.Timeout)
	defer cancel()

	filter := bson.M{
		"key": s.PosKey,
	}

	set := bson.M{
		"key":   s.PosKey,
		"value": value,
	}

	set["version"] = time.Now().Unix()
	if len(unix) > 0 {
		set["version"] = unix[0]
	}

	update := bson.M{
		"$set": set,
	}

	upsert := true
	updOpt := &options.UpdateOptions{Upsert: &upsert}
	if _, err := s.coll.UpdateOne(ctx, filter, update, updOpt); err != nil {
		Logger.Error("set pos str err", zap.Error(err))
		return err
	}
	return nil
}

func (s *SyncerPosMgoStore) Pos() (*SyncerPos, error) {
	syncTimeStr, version, err := s.getPosStr()
	if err != nil {
		return nil, err
	}

	Logger.Info("query latest sync time ===> ",
		zap.String("syncTimeStr", syncTimeStr),
		zap.String("key", s.PosKey))

	syncTimeStrs := strings.Split(syncTimeStr, ":")
	if len(syncTimeStr) != 3 {
		Logger.Error("syncTimeStr invalid",
			zap.Strings("syncTimeStrs", syncTimeStrs))
		return nil, errors.New("syncTimeStr invalid")
	}

	pos, err := strconv.ParseInt(syncTimeStrs[2], 10, 64)
	if err != nil {
		Logger.Warn("pos parse err", zap.Error(err))
		return nil, err
	}
	syncStart, err := strconv.ParseInt(syncTimeStrs[0], 10, 64)
	if err != nil {
		Logger.Warn("syncTimeStr parse err", zap.Error(err))
		return nil, err
	}
	syncEnd, err := strconv.ParseInt(syncTimeStrs[1], 10, 64)
	if err != nil {
		Logger.Warn("syncTimeStr parse err", zap.Error(err))
		return nil, err
	}
	return &SyncerPos{
		SyncStartTimestamp: syncStart,
		SyncEndTimestamp:   syncEnd,
		Pos:                pos,
		LastSyncTime:       version,
	}, nil
}

func (s *SyncerPosMgoStore) Save(pos interface{}, unix ...int64) error {
	return s.setPosStr(pos, unix...)
}
