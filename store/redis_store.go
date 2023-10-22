package store

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"

	. "github.com/Rayson9418/cdc/logger"
	cdcredis "github.com/Rayson9418/cdc/redis"
)

type BinlogRedisStore struct {
	logPosKey string
	rdb       *redis.Client
}

func NewBinlogRedisStore(key string) *BinlogRedisStore {
	return &BinlogRedisStore{
		logPosKey: key,
		rdb:       cdcredis.GetRedisClient(),
	}
}

func (b *BinlogRedisStore) Pos() (string, uint32, error) {
	ctx := context.Background()

	binlog, err := b.rdb.Get(ctx, b.logPosKey).Result()
	if err != nil {
		if err == redis.Nil {
			return "", 0, nil
		}
		Logger.Warn("binlog pos key not found",
			zap.String("pos_key", b.logPosKey),
			zap.Error(err))
		return "", 0, err
	}

	binlogArr := strings.Split(binlog, ":")
	if 2 != len(binlogArr) {
		Logger.Warn("split binlog pos value err",
			zap.Error(err))
		return "", 0, errors.New("parse binlog err")
	}

	binlogFile := binlogArr[0]
	binlogPosInt, err := strconv.Atoi(binlogArr[1])
	if err != nil {
		Logger.Warn("parse bin log position err", zap.Error(err))
		return "", 0, err
	}

	return binlogFile, uint32(binlogPosInt), nil
}

func (b *BinlogRedisStore) Save(file string, pos uint32) error {
	return b.rdb.Set(context.Background(), b.logPosKey,
		fmt.Sprintf("%s:%d", file, pos), 0).Err()
}

type StreamRedisStore struct {
	logPosKey string
	rdb       *redis.Client
}

func NewStreamRedisStore(key string) *StreamRedisStore {
	return &StreamRedisStore{
		logPosKey: key,
		rdb:       cdcredis.GetRedisClient(),
	}
}

func (s *StreamRedisStore) Pos() (string, error) {
	ctx := context.Background()

	token, err := s.rdb.Get(ctx, s.logPosKey).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil
		}
		Logger.Warn("resume tokem key not found",
			zap.String("token_key", s.logPosKey),
			zap.Error(err))
		return "", err
	}

	return token, nil
}

func (s *StreamRedisStore) Save(token string) error {
	return s.rdb.Set(context.Background(), s.logPosKey,
		token, 0).Err()
}
