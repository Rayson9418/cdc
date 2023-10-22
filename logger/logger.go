package logger

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Logger *zap.Logger
	once   sync.Once
)

func init() {
	once.Do(func() {
		InitLogger()
	})
}

func InitLogger() {
	encoderCofig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	ws := zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout))

	// set log level
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(defaultLogLevel())

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCofig), // encoder config
		ws,                                   // writeSyncer
		atomicLevel,                          // log level
	)

	// show file name and line number
	caller := zap.AddCaller()
	// enable development mode, which shows the caller in the log
	development := zap.Development()
	// add default app name to log
	//field := zap.Fields(zap.String("app", appName))
	Logger = zap.New(
		core,
		caller,
		development,
	)
}

func SetLogger(logger *zap.Logger) {
	Logger = logger
}

func defaultLogLevel() zapcore.Level {
	return zap.DebugLevel
}
