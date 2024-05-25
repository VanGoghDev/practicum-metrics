package app

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/VanGoghDev/practicum-metrics/internal/server/config"
	"github.com/VanGoghDev/practicum-metrics/internal/server/services/filesaver"
	"github.com/VanGoghDev/practicum-metrics/internal/storage/memstorage"
)

type FileSaver interface {
	Run() error
}

type App struct {
	log       *zap.Logger
	filesaver FileSaver
}

func New(cfg *config.Config, log *zap.Logger, mstrg *memstorage.MemStorage) (*App, error) {
	// init services.
	filesvr, err := filesaver.New(cfg, log, mstrg)
	if err != nil {
		return nil, fmt.Errorf("failed to init filesaver service: %w", err)
	}

	return &App{
		filesaver: filesvr,
		log:       log,
	}, nil
}

func (a *App) RunApp() error {
	err := a.filesaver.Run()
	if err != nil {
		return fmt.Errorf("App.Run: %w", err)
	}
	return nil
}
