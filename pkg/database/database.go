package database

import (
	"database/sql"
	"enterprisedata-exchange/internal/config"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func Connect(cfg *config.Config) (*sql.DB, error) {
	const op = "pkg.database.Connect"
	db, err := sql.Open("sqlite3", cfg.StoragePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return db, nil
}
