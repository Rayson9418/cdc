package redis

import (
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"

	. "github.com/Rayson9418/cdc/logger"
)

const (
	KMasterName        = "mymaster"
	KRedisTypeSingle   = "single"
	KRedisTypeSentinel = "sentinel"
)

var (
	initOnce       sync.Once
	globalRedisCli *redis.Client
)

func NewDefaultOpt() *Options {
	opt := new(Options)

	opt.Type = KRedisTypeSingle
	opt.Addr = "127.0.0.1:2369"
	opt.Pwd = "123456"
	opt.Auth = true

	return opt
}

func InitClient(opt *Options) error {
	var (
		err    error
		client *redis.Client
	)

	initOnce.Do(func() {
		switch opt.Type {
		case KRedisTypeSingle:
			client, err = initSingle(opt)
		case KRedisTypeSentinel:
			client, err = initSentinels(opt)
		default:
			err = errors.New("unknown type of redis")
		}
		if err != nil {
			return
		}
		globalRedisCli = client
	})
	return err
}

func initSingle(opt *Options) (*redis.Client, error) {
	connOpt := &redis.Options{
		Addr:     opt.Addr,
		Password: opt.Pwd,
	}
	if !opt.Auth {
		connOpt.Password = ""
	}

	redisCli := redis.NewClient(connOpt)
	if err := redisCli.Ping(redisCli.Context()).Err(); err != nil {
		Logger.Error("failed to new redis client", zap.Error(err))
		return nil, err
	}

	return redisCli, nil
}

func initSentinels(opt *Options) (*redis.Client, error) {
	failOverOptions := &redis.FailoverOptions{
		MasterName:       KMasterName,
		SentinelAddrs:    strings.Split(opt.Addr, ","),
		Password:         opt.Pwd,
		SentinelPassword: opt.Pwd,
	}

	redisCli := redis.NewFailoverClient(failOverOptions)
	if err := redisCli.Ping(redisCli.Context()).Err(); err != nil {
		Logger.Error("failed to new redis client", zap.Error(err))
		return nil, err
	}

	return redisCli, nil
}

func GetRedisClient() *redis.Client {
	return globalRedisCli
}

func SetNx(key string, value interface{}, seconds time.Duration) (bool, error) {
	return globalRedisCli.SetNX(globalRedisCli.Context(), key, value, seconds).Result()
}

func Del(keys ...string) error {
	return globalRedisCli.Del(globalRedisCli.Context(), keys...).Err()
}
