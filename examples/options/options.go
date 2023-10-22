package options

import (
	"os"

	. "github.com/Rayson9418/cdc/logger"
	cdcmongo "github.com/Rayson9418/cdc/mongo"
	cdcmysql "github.com/Rayson9418/cdc/mysql"
	cdcredis "github.com/Rayson9418/cdc/redis"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

type CdcOptions struct {
	Mysql *cdcmysql.Options `yaml:"mysql"`
	Mongo *cdcmongo.Options `yaml:"mongo"`
	Redis *cdcredis.Options `yaml:"redis"`
}

var CdcOpt = &CdcOptions{}

func Init() error {
	data, err := os.ReadFile("cdc.yaml")
	if err != nil {
		Logger.Error("read file err", zap.Error(err))
		return err
	}

	if err = yaml.Unmarshal(data, &CdcOpt); err != nil {
		Logger.Error("yaml unmarshal err", zap.Error(err))
		return err
	}
	return nil
}
