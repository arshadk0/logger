package network_logger

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	NETWORK_LOGGER_TYPE     = "network"
	RESPONSE_MESSAGE_KEY    = "log-message-key"
	RESPONSE_ERROR_CODE_KEY = "log-error-code-key"
	LEVEL_KEY               = "level"
	MESSAGE_KEY             = "type"
	TIMESTAMP_KEY           = "timestamp"
)

var networkLogger *zap.Logger

func init() {
	conf := zap.Config{
		Level:             zap.NewAtomicLevelAt(zapcore.InfoLevel),
		Encoding:          "json",
		DisableCaller:     true,
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
	networkLogger = zap.Must(conf.Build())
}

func NetworkLogger(c *gin.Context) {
	startTime := time.Now()
	var fields []zapcore.Field
	fields = append(fields, zap.String("endpoint", c.Request.URL.String()))
	fields = append(fields, zap.String("method", c.Request.Method))
	fields = append(fields, zap.String("ip", c.ClientIP()))
	fields = append(fields, zap.String("user-agent", c.Request.UserAgent()))
	c.Next()
	fields = append(fields, zap.Int("status", c.Writer.Status()))
	message, exists := c.Get(RESPONSE_MESSAGE_KEY)
	if exists {
		fields = append(fields, zap.String("message", message.(string)))
	}
	errorCode, exists := c.Get(RESPONSE_ERROR_CODE_KEY)
	if exists {
		fields = append(fields, zap.String("error-code", errorCode.(string)))
	}
	elapsed := time.Since(startTime)
	fields = append(fields, zap.Float64("response-time-ms", float64(elapsed.Nanoseconds())/1e6))
	fields = append(fields, zap.Int64("response-size-bytes", int64(c.Writer.Size())))
	networkLogger.Info(NETWORK_LOGGER_TYPE, fields...)
}
