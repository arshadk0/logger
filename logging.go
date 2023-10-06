package logging

import (
	"context"
	"errors"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

const (
	DATABASE_LOGGER_TYPE = "db"
)

type GormLogger struct {
	ZapLogger                 *zap.Logger
	LogLevel                  zapcore.Level
	SlowThreshold             time.Duration
	IgnoreRecordNotFoundError bool
}

var DbLogger GormLogger

func init() {
	conf := zap.Config{
		Level:             zap.NewAtomicLevelAt(zapcore.InfoLevel),
		Encoding:          "json",
		DisableCaller:     false,
		DisableStacktrace: false,
		OutputPaths:       []string{"stdout"},
		ErrorOutputPaths:  []string{"stderr"},
		Development:       false,
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:     TIMESTAMP_KEY,
			EncodeTime:  zapcore.RFC3339TimeEncoder,
			LevelKey:    LEVEL_KEY,
			EncodeLevel: zapcore.CapitalLevelEncoder,
			MessageKey:  MESSAGE_KEY,
		},
	}
	zapLogger := zap.Must(conf.Build())
	DbLogger = GormLogger{
		ZapLogger:                 zapLogger,
		SlowThreshold:             200 * time.Millisecond,
		IgnoreRecordNotFoundError: false,
	}
}

func (l GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
}

func (l GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
}

func (l GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
}

func (l GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	switch {
	case level == logger.Error:
		l.LogLevel = zapcore.ErrorLevel
	case level == logger.Warn:
		l.LogLevel = zapcore.WarnLevel
	case level == logger.Info:
		l.LogLevel = zapcore.InfoLevel
	default:
		l.LogLevel = zapcore.DebugLevel
	}
	return l
}

func (l GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= zapcore.DebugLevel {
		return
	}
	sql, rows := fc()
	elapsed := time.Since(begin)
	switch {
	case err != nil && l.LogLevel <= zapcore.ErrorLevel && (!errors.Is(err, gorm.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		if rows == -1 {
			l.ZapLogger.Error(DATABASE_LOGGER_TYPE, zap.Error(err), zap.String("sql", sql), zap.Float64("elapsed-in-ms", float64(elapsed.Nanoseconds())/1e6), zap.String("rows", "-"), zap.String("caller", utils.FileWithLineNum()))
		} else {
			l.ZapLogger.Error(DATABASE_LOGGER_TYPE, zap.Error(err), zap.String("sql", sql), zap.Float64("elapsed-in-ms", float64(elapsed.Nanoseconds())/1e6), zap.Int64("rows", rows), zap.String("caller", utils.FileWithLineNum()))
		}
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel <= zapcore.WarnLevel:
		if rows == -1 {
			l.ZapLogger.Warn(DATABASE_LOGGER_TYPE, zap.String("sql", sql), zap.Float64("elapsed-in-ms", float64(elapsed.Nanoseconds())/1e6), zap.String("rows", "-"), zap.String("caller", utils.FileWithLineNum()))
		} else {
			l.ZapLogger.Warn(DATABASE_LOGGER_TYPE, zap.String("sql", sql), zap.Float64("elapsed-in-ms", float64(elapsed.Nanoseconds())/1e6), zap.Int64("rows", rows), zap.String("caller", utils.FileWithLineNum()))
		}
	case l.LogLevel <= zapcore.InfoLevel:
		if rows == -1 {
			l.ZapLogger.Info(DATABASE_LOGGER_TYPE, zap.String("sql", sql), zap.Float64("elapsed-in-ms", float64(elapsed.Nanoseconds())/1e6), zap.String("rows", "-"), zap.String("caller", utils.FileWithLineNum()))
		} else {
			l.ZapLogger.Info(DATABASE_LOGGER_TYPE, zap.String("sql", sql), zap.Float64("elapsed-in-ms", float64(elapsed.Nanoseconds())/1e6), zap.Int64("rows", rows), zap.String("caller", utils.FileWithLineNum()))
		}
	}
}
