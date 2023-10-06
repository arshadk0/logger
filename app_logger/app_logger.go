package app_logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	LEVEL_KEY               = "level"
	MESSAGE_KEY             = "type"
	TIMESTAMP_KEY           = "timestamp"
	APPLICATION_LOGGER_TYPE = "app"
)

var ApplicationLogger *zap.Logger

func init() {
	conf := zap.Config{
		Level:             zap.NewAtomicLevelAt(zapcore.InfoLevel),
		Encoding:          "json",
		DisableCaller:     true,
		DisableStacktrace: false,
		Development:       false,
		OutputPaths:       []string{"stdout"},
		ErrorOutputPaths:  []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:     TIMESTAMP_KEY,
			EncodeTime:  zapcore.RFC3339TimeEncoder,
			LevelKey:    LEVEL_KEY,
			EncodeLevel: zapcore.CapitalLevelEncoder,
			MessageKey:  MESSAGE_KEY,
		},
	}
	ApplicationLogger = zap.Must(conf.Build())
}

func Info(info string) {
	ApplicationLogger.Info(APPLICATION_LOGGER_TYPE, zap.String("message", info))
}

func Warn(info string) {
	ApplicationLogger.Warn(APPLICATION_LOGGER_TYPE, zap.String("message", info))
}

func Error(msg string, err error) {
	ApplicationLogger.Error(APPLICATION_LOGGER_TYPE, zap.String("message", msg), zap.Error(err))
}

func Fatal(msg string, err error) {
	ApplicationLogger.Fatal(APPLICATION_LOGGER_TYPE, zap.String("message", msg), zap.Error(err))
}
