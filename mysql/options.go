package mysql

import (
	"github.com/go-mysql-org/go-mysql/schema"
)

type Options struct {
	Addr      string      `yaml:"addr"`
	User      string      `yaml:"user"`
	Pwd       string      `yaml:"pwd"`
	Charset   string      `yaml:"charset"`
	Databases []*Database `yaml:"databases"`
	//MonitorTableMap map[string]map[string]*TableInfo // table_name -> action -> table_info
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

//var opt = &Options{}

func NewDefaultOpt() *Options {
	opt := new(Options)

	opt.Addr = "127.0.0.1:27017"
	opt.User = "root"
	opt.Pwd = "123456"
	return opt
}
