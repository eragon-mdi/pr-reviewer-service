package logger

import (
	"github.com/eragon-mdi/go-playground/logging"
	"github.com/eragon-mdi/pr-reviewer-service/internal/common/configs"
	"github.com/go-faster/errors"
	"go.uber.org/zap"
)

const ErrInitLogger = "failed to init logger"

func New(cfg configs.Logger) (*zap.SugaredLogger, error) {
	l, err := logging.NewLogger(cfg.Level, cfg.Encoding, cfg.Output, cfg.MessageKey)
	if err != nil {
		return nil, errors.Wrap(err, ErrInitLogger)
	}

	return l, nil
}
