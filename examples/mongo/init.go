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

	baseopt "examples/options"
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

type Option struct {
	Addr     string `yaml:"addr"`
	User     string `yaml:"user"`
	Pwd      string `yaml:"pwd"`
	Auth     bool   `yaml:"auth"`
	Direct   bool   `yaml:"direct"`
	PoolSize uint64 `yaml:"pool_size"`
	Timeout  uint64 `yaml:"timeout"`
	uri      string
}

func GetTimeoutLimit() time.Duration {
	return timeout
}

func initClient(opt *Option) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cli, err := mongo.Connect(ctx, options.Client().ApplyURI(opt.uri).SetMaxPoolSize(opt.PoolSize).SetDirect(opt.Direct))
	if err != nil {
		Logger.Warn("connect mongodb err", zap.Error(err))
		return nil, err
	}
	return cli, nil
}

func InitClient() error {
	opt, err := getMongoOpt()
	if err != nil {
		Logger.Error("get mysql options err", zap.Error(err))
		return err
	}

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

func getMongoOpt() (*Option, error) {
	opt := &Option{
		Addr:     baseopt.CdcOpt.Mongo.Addr,
		User:     baseopt.CdcOpt.Mongo.User,
		Pwd:      baseopt.CdcOpt.Mongo.Pwd,
		Auth:     baseopt.CdcOpt.Mongo.Auth,
		Direct:   baseopt.CdcOpt.Mongo.Direct,
		PoolSize: baseopt.CdcOpt.Mongo.PoolSize,
		Timeout:  baseopt.CdcOpt.Mongo.Timeout,
	}
	return opt, nil
}

func GetClient() *mongo.Client {
	return globalClient
}

func SetClient(cli *mongo.Client) {
	globalClient = cli
}
