package shared

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
	"go.uber.org/fx"
)

type migrationStorageImpl struct {
	db     *sql.DB
	logger *log.Logger
}

func NewMigrationStorage(cfg *Config, logger *log.Logger, lc fx.Lifecycle) (MigrationStorage, error) {
	db, err := sql.Open("sqlite3", cfg.SqlitePath)
	if err != nil {
		return nil, err
	}

	// Ensure migrations table exists
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS migrations (
			identifier TEXT PRIMARY KEY,
			applied INTEGER
		)
	`)
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return db.Close()
		},
	})

	return &migrationStorageImpl{db: db, logger: logger}, nil
}

func (ms *migrationStorageImpl) LoadState() (map[string]bool, error) {
	rows, err := ms.db.Query("SELECT identifier, applied FROM migrations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	state := make(map[string]bool)
	for rows.Next() {
		var identifier string
		var applied int
		if err := rows.Scan(&identifier, &applied); err != nil {
			return nil, err
		}
		state[identifier] = applied == 1
	}

	return state, nil
}

func (ms *migrationStorageImpl) SaveState(identifier string, applied bool) error {
	_, err := ms.db.Exec(`
		INSERT INTO migrations (identifier, applied)
		VALUES (?, ?)
		ON CONFLICT(identifier) DO UPDATE SET applied = ?
	`, identifier, boolToInt(applied), boolToInt(applied))
	return err
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
