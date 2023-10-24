package mongo

import "fmt"

type Options struct {
	Addr      string      `yaml:"addr"`
	User      string      `yaml:"user"`
	Pwd       string      `yaml:"pwd"`
	Auth      bool        `yaml:"auth"`
	Direct    bool        `yaml:"direct"`
	PoolSize  uint64      `yaml:"pool_size"`
	Timeout   uint64      `yaml:"timeout"`
	Databases []*Database `yaml:"databases"`
	uri       string
}

type Database struct {
	Name        string            `yaml:"name"`
	Collections []*CollectionInfo `yaml:"collections"`
}

type CollectionInfo struct {
	Name    string   `yaml:"name"`
	Actions []string `yaml:"actions"`
}

func NewDefaultOpt() *Options {
	opt := new(Options)

	opt.Addr = "127.0.0.1:27017"
	opt.User = "root"
	opt.Pwd = "123456"
	opt.Auth = true
	opt.PoolSize = 100
	opt.Timeout = 60

	opt.uri = fmt.Sprintf(kUriFmt, opt.User, opt.Pwd, opt.Addr)
	if !opt.Auth {
		opt.uri = fmt.Sprintf(kNoAuthUriFmt, opt.Addr)
	}
	return opt
}
