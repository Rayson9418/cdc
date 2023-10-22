package syncer

import (
	"sync"

	"github.com/spf13/viper"
	"go.uber.org/zap"

	. "github.com/Rayson9418/cdc/logger"
)

var (
	opt  *options
	once sync.Once
)

type options struct {
	StartHour  int `mapstructure:"start_hour"`
	EndHour    int `mapstructure:"end_hour"`
	BatchLimit int `mapstructure:"batch_limit"`
	Interval   int `mapstructure:"interval"`
}

func initOnce() error {
	vp := viper.New()
	vp.SetConfigType("yaml")
	vp.SetConfigFile("/yzp/base/env.yaml")
	err := vp.ReadInConfig()
	if err != nil {
		Logger.Warn("read config err", zap.Error(err))
		return err
	}

	opt = &options{}
	if err = vp.UnmarshalKey("sinker", opt); err != nil {
		Logger.Warn("unmarshal sinker conf err", zap.Error(err))
		return err
	}
	return nil
}

func initOpt() error {
	var err error
	once.Do(func() {
		err = initOnce()
	})
	return err
}
