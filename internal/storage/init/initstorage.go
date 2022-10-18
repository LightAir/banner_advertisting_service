package initstorage

import (
	"fmt"

	"github.com/LightAir/bas/internal/config"
	"github.com/LightAir/bas/internal/core"
	"github.com/LightAir/bas/internal/logger"
	memorystorage "github.com/LightAir/bas/internal/storage/memory"
	sqlstorage "github.com/LightAir/bas/internal/storage/sql"
)

func NewStorage(cfg *config.Config, logg *logger.Logger) (core.Storage, error) {
	switch cfg.DB.Type {
	case "mem":
		return memorystorage.New(logg), nil
	case "sql":
		return sqlstorage.New(cfg, config.GetDsn(cfg.DB.SQL)), nil
	}

	return nil, fmt.Errorf("unknown database type: %q", cfg.DB.Type)
}
