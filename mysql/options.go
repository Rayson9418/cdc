package mysql

import (
	"github.com/go-mysql-org/go-mysql/schema"
)

type Options struct {
	Addr       string      `yaml:"addr"`
	User       string      `yaml:"user"`
	Pwd        string      `yaml:"pwd"`
	Charset    string      `yaml:"charset"`
	MaxIdleNum int         `yaml:"max_idle_num"`
	MaxConnNum int         `yaml:"max_conn_num"`
	DefaultDb  string      `yaml:"default_db"`
	Databases  []*Database `yaml:"databases"`
}

type Database struct {
	Name   string       `yaml:"name"`
	Tables []*TableInfo `yaml:"tables"`
}

type TableInfo struct {
	Name          string   `yaml:"name"`
	Actions       []string `yaml:"actions"`
	ColumnMapping map[string]*columnInfo
}

type columnInfo struct {
	ColumnIndex    int
	ColumnMetadata *schema.TableColumn
}

func NewDefaultOpt() *Options {
	opt := new(Options)

	opt.Addr = "127.0.0.1:27017"
	opt.User = "root"
	opt.Pwd = "123456"
	opt.MaxIdleNum = 100
	opt.MaxConnNum = 1000
	opt.DefaultDb = "demo"

	return opt
}
