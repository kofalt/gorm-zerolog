package gorm_zerolog

import (
	"context"
	"errors"
	"time"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

type Logger struct {
	SlowThreshold         time.Duration
	SourceField           string
	SkipErrRecordNotFound bool
	Logger                zerolog.Logger
}

func New(l zerolog.Logger) *Logger {
	return &Logger{
		Logger: l,
	}
}

func (l *Logger) LogMode(gormlogger.LogLevel) gormlogger.Interface {
	return l
}

func (l *Logger) Info(ctx context.Context, s string, args ...interface{}) {
	l.Logger.Info().Msgf(s, args)
}

func (l *Logger) Warn(ctx context.Context, s string, args ...interface{}) {
	l.Logger.Warn().Msgf(s, args)
}

func (l *Logger) Error(ctx context.Context, s string, args ...interface{}) {
	l.Logger.Error().Msgf(s, args)
}

func (l *Logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, _ := fc()

	fields := map[string]interface{}{
		"elapsed": elapsed,
		"sql":     sql,
	}

	if l.SourceField != "" {
		fields[l.SourceField] = utils.FileWithLineNum()
	}

	if err != nil && !(l.SkipErrRecordNotFound && errors.Is(err, gorm.ErrRecordNotFound)) {
		l.Logger.Error().Err(err).Fields(fields).Msg("SQL error")
		return
	}

	if l.SlowThreshold != 0 && elapsed > l.SlowThreshold {
		l.Logger.Warn().Fields(fields).Msgf("SQL slow")
		return
	}

	l.Logger.Trace().Fields(fields).Msgf("SQL")
}
