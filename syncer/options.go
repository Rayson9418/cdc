package syncer

var (
	opt *Options
)

type Options struct {
	StartHour  int `yaml:"start_hour"`
	EndHour    int `yaml:"end_hour"`
	BatchLimit int `yaml:"batch_limit"`
	Interval   int `yaml:"interval"`
}

func NewDefaultOpt() *Options {
	opt = new(Options)

	opt.StartHour = 0
	opt.EndHour = 5
	opt.BatchLimit = 2
	opt.Interval = 30

	return opt
}
