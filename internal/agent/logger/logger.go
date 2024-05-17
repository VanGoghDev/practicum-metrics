package logger

import (
	"fmt"

	"go.uber.org/zap"
)

func New(level string) (*zap.Logger, error) {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return nil, fmt.Errorf("failed to parse given level = %s: %w", level, err)
	}

	cfg := zap.NewProductionConfig()
	cfg.Level = lvl

	zl, err := cfg.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build log cfg: %w", err)
	}

	defer func() {
		errS := zl.Sync()
		if err == nil && errS != nil {
			err = fmt.Errorf("failed to sync log: %w", errS)
		}
	}()

	return zl, err
}
