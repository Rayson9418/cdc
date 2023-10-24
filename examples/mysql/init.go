package mysql

import (
	"context"
	"fmt"
	"sync"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	. "examples/logger"
	"examples/options"
)

const (
	KMysqlUrlTemplate = "%s:%s@tcp(%s)/%s"
	KMysqlUrlSuffix   = "?charset=utf8&parseTime=true&loc=Asia%2FShanghai"
)

var (
	initOnce sync.Once
	globalDB *gorm.DB
)

type Option struct {
	Addr      string // 地址
	Username  string // 用户名
	Password  string // 密码
	DefaultDB string // 默认数据Once
	MaxIdle   int    // 最大空闲时间
	MaxConn   int    // 最大连接数
}

func InitClient() error {
	var err error
	initOnce.Do(func() {
		globalDB, err = initClient()
	})

	return err
}

func initClient() (*gorm.DB, error) {
	opt, err := GetMysqlOpt()
	if err != nil {
		Logger.Error("get mysql options err", zap.Error(err))
		return nil, err
	}

	url := fmt.Sprintf(KMysqlUrlTemplate, opt.Username,
		opt.Password, opt.Addr, opt.DefaultDB)
	url += KMysqlUrlSuffix
	globalDB, err = gorm.Open(mysql.Open(url), &gorm.Config{})
	if err != nil {
		Logger.Error("connect mysql err", zap.Error(err))
		return nil, err
	}

	db, err := globalDB.DB()
	if err != nil {
		Logger.Error("failed to set db", zap.String("err", err.Error()))
		return nil, err
	}

	// Set up the connection pool
	db.SetMaxOpenConns(opt.MaxConn)
	db.SetMaxIdleConns(opt.MaxIdle)
	return globalDB, nil
}

func GetMysqlOpt() (*Option, error) {
	opt := &Option{
		Addr:      options.CdcOpt.Mysql.Addr,
		Username:  options.CdcOpt.Mysql.User,
		Password:  options.CdcOpt.Mysql.Pwd,
		DefaultDB: "demo",
		MaxIdle:   100,
		MaxConn:   1000,
	}
	return opt, nil
}

func GetDB(ctx context.Context) *gorm.DB {
	return globalDB.WithContext(ctx)
}
