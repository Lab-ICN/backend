package logging

import (
	"github.com/Lab-ICN/backend/user-service/internal/config"
	"go.uber.org/zap"
)

func New(cfg *config.Config) (*zap.Logger, error) {
	level, err := zap.ParseAtomicLevel(cfg.Logging.Level)
	if err != nil {
		return nil, err
	}
	_cfg := zap.Config{
		Level:             level,
		Development:       cfg.Logging.Development,
		DisableCaller:     cfg.Logging.DisableCaller,
		DisableStacktrace: cfg.Logging.DisableStacktrace,
		Encoding:          cfg.Logging.Encoding,
		OutputPaths:       cfg.Logging.OutputPaths,
		ErrorOutputPaths:  cfg.Logging.ErrorOutputPaths,
	}
	return _cfg.Build()
}
