package redis

const (
	kDefaultRedisPoolSize           = 100
	kDefaultRedisIdleTimeout        = 600
	kDefaultRedisDialTimeout        = 10
	kDefaultRedisReadTimeout        = 10
	kDefaultRedisWriteTimeout       = 10
	kDefaultRedisIdleCheckFrequency = 60
	kDefaultRedisConnMaxRetries     = 3
	kDefaultRedisDB                 = 0
	kRedisKeyExist                  = 1
)

type Options struct {
	Type      string `yaml:"type"`
	Addr      string `yaml:"addr"`
	Pwd       string `yaml:"pwd"`
	Auth      bool   `yaml:"auth"`
	TLSConfig bool   `yaml:"tlsconfig"`
}
