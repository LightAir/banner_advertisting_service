package sqlstorage

import (
	"context"
	"fmt"

	"github.com/LightAir/bas/internal/config"
	"github.com/jmoiron/sqlx"
)

type Storage struct {
	dsn    string
	db     *sqlx.DB
	config config.Config
	ctx    context.Context
}

func (s *Storage) Connect(ctx context.Context) error {
	db, err := sqlx.Open("pgx", s.dsn)
	if err != nil {
		return fmt.Errorf("connection error: %w", err)
	}

	s.db = db
	s.ctx = ctx

	return nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func New(config *config.Config, dsn string) *Storage {
	return &Storage{
		dsn:    dsn,
		config: *config,
		db:     nil,
	}
}
