package options

import (
	"os"

	cdcmongo "github.com/Rayson9418/cdc/mongo"
	cdcmysql "github.com/Rayson9418/cdc/mysql"
	cdcredis "github.com/Rayson9418/cdc/redis"
	cdcsyncer "github.com/Rayson9418/cdc/syncer"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"

	. "examples/logger"
)

type CdcOptions struct {
	Mysql  *cdcmysql.Options  `yaml:"mysql"`
	Mongo  *cdcmongo.Options  `yaml:"mongo"`
	Redis  *cdcredis.Options  `yaml:"redis"`
	Syncer *cdcsyncer.Options `yaml:"syncer"`
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
