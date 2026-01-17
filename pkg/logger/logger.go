package logger

import (
	"go.uber.org/zap"
)

var (
	Log   *zap.Logger = zap.NewNop()
	Sugar *zap.SugaredLogger
)

// New - иницализация логгера.
func New(level string) error {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}
	cfg := zap.NewProductionConfig()
	cfg.Level = lvl
	zl, err := cfg.Build()
	if err != nil {
		return err
	}
	Log = zl
	Sugar = Log.Sugar()
	return nil
}
