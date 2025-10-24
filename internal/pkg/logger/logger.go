package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	gorm_zerolog "github.com/wei840222/gorm-zerolog"
	gormlogger "gorm.io/gorm/logger"
)

type Logger struct {
	Level zerolog.Level
}

type LoggerInterface interface {
	Setup()
	GetLogger() zerolog.Logger
	GetDatabaseLogger() gormlogger.Interface
}

func NewLogger(level string) *Logger {
	return &Logger{
		Level: getLevel(level),
	}
}

func (l *Logger) Setup() {
	zerolog.SetGlobalLevel(l.Level)

	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	})
}

func (l *Logger) GetLogger() zerolog.Logger {
	return log.Logger
}

func (l *Logger) GetDatabaseLogger() gormlogger.Interface {
	return gorm_zerolog.New()
}

func getLevel(level string) zerolog.Level {
	switch level {
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "debug":
		return zerolog.DebugLevel
	case "trace":
		return zerolog.TraceLevel
	default:
		return zerolog.InfoLevel
	}
}
