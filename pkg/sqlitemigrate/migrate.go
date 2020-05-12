package sqlitemigrate

import (
	"context"
	"database/sql"
	"fmt"
)

type MigrationRegistry struct {
	migrations []Migration
}

func NewRegistry() *MigrationRegistry {
	return &MigrationRegistry{
		migrations: make([]Migration, 0, 10),
	}
}

func (r *MigrationRegistry) GetCurrentVersion(ctx context.Context, db *sql.DB) (int, error) {
	var version int
	if err := db.QueryRowContext(ctx, "PRAGMA user_version").Scan(&version); err != nil {
		return -1, err
	}
	return version, nil
}

func (r *MigrationRegistry) Apply(ctx context.Context, db *sql.DB) error {
	current, err := r.GetCurrentVersion(ctx, db)
	if err != nil {
		return err
	}
	for idx, mig := range r.migrations {
		mig.Version = idx + 1
		if current >= idx+1 {
			continue
		}
		if err := mig.Apply(ctx, db); err != nil {
			return err
		}
	}
	return nil
}

func (r *MigrationRegistry) RegisterMigration(ups []string, downs []string) {
	m := Migration{
		Up:   ups,
		Down: downs,
	}
	r.migrations = append(r.migrations, m)
}

type Migration struct {
	Version int
	Up      []string
	Down    []string
}

func (m *Migration) Apply(ctx context.Context, db *sql.DB) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	for _, up := range m.Up {
		if _, err := tx.ExecContext(ctx, up); err != nil {
			tx.Rollback()
			return err
		}
	}
	if _, err := tx.ExecContext(ctx, fmt.Sprintf("PRAGMA user_version = %d", m.Version)); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}
