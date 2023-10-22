package mongo

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"

	. "github.com/Rayson9418/cdc/logger"
)

const (
	kUriFmt       = "mongodb://%s:%s@%s"
	kNoAuthUriFmt = "mongodb://%s"
)

var (
	initOnce     sync.Once
	globalClient *mongo.Client
	timeout      time.Duration
)

func GetTimeoutLimit() time.Duration {
	return timeout
}

func initClient(opt *Options) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cli, err := mongo.Connect(ctx, options.Client().ApplyURI(opt.uri).SetMaxPoolSize(opt.PoolSize).SetDirect(opt.Direct))
	if err != nil {
		Logger.Warn("connect mongodb err", zap.Error(err))
		return nil, err
	}
	return cli, nil
}

func InitClient(opt *Options) error {
	var err error
	initOnce.Do(func() {
		opt.uri = fmt.Sprintf(kUriFmt, opt.User, opt.Pwd, opt.Addr)
		if !opt.Auth {
			opt.uri = fmt.Sprintf(kNoAuthUriFmt, opt.Addr)
		}

		timeout = time.Duration(opt.Timeout) * time.Second
		globalClient, err = initClient(opt)
	})
	return err
}

func GetClient() *mongo.Client {
	return globalClient
}

func SetClient(cli *mongo.Client) {
	globalClient = cli
}
