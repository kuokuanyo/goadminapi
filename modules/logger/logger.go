package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	defaultEncoderCfg = EncoderCfg{
		TimeKey:       "ts",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		Level:         "capitalColor",
		Time:          "ISO8601",
		Duration:      "seconds",
		Caller:        "short",
		Encoding:      "console",
	}
	defaultRotateCfg = RotateCfg{
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   false,
	}
	logger = &Logger{
		rotate:  defaultRotateCfg,
		encoder: defaultEncoderCfg,
		Level:   zapcore.InfoLevel,
	}
)

type EncoderCfg struct {
	TimeKey       string
	LevelKey      string
	NameKey       string
	CallerKey     string
	MessageKey    string
	StacktraceKey string
	Level         string
	Time          string
	Duration      string
	Caller        string
	Encoding      string
}

type RotateCfg struct {
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
}

type Logger struct {
	logger        *zap.Logger
	sugaredLogger *zap.SugaredLogger

	infoLogOff   bool
	errorLogOff  bool
	accessLogOff bool

	accessAssetsLogOff bool

	debug bool

	sqlLogOpen bool

	infoLogPath   string
	errorLogPath  string
	accessLogPath string

	rotate  RotateCfg
	encoder EncoderCfg

	Level zapcore.Level
}

// Error print the error message.
func Error(err ...interface{}) {
	if !logger.errorLogOff && logger.Level <= zapcore.ErrorLevel {
		logger.sugaredLogger.Error(err...)
	}
}
