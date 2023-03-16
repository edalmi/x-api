package server

import (
	"log"

	"github.com/edalmi/x-api/config"
	"github.com/edalmi/x-api/logging"
	stdlog "github.com/edalmi/x-api/logging/log"
	_slog "github.com/edalmi/x-api/logging/slog"
	zaplog "github.com/edalmi/x-api/logging/zap"
	"go.uber.org/zap"
	"golang.org/x/exp/slog"
)

func setupLogger(mode string, cfg *config.Logger) (logging.Logger, error) {
	if cfg.Std != nil {
		return stdlog.New(log.Default()), nil
	}

	if cfg.Slog != nil {
		return _slog.New(slog.Default()), nil
	}

	if cfg.Zap != nil {
		logger, err := zap.NewProduction()
		if err != nil {
			return nil, err
		}

		if mode == config.ModeDev {
			logger, err = zap.NewDevelopment()
			if err != nil {
				return nil, err
			}
		}

		return zaplog.New(logger), nil
	}

	return stdlog.New(log.Default()), nil
}
