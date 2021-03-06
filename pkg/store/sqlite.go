package store

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/zerok/geotrace/pkg/sqlitemigrate"
)

var mig *sqlitemigrate.MigrationRegistry

type sqliteStore struct {
	path string
	db   *sql.DB
}

func NewSQLiteStore(path string) Store {
	return &sqliteStore{path: path}
}

func (s *sqliteStore) Add(t time.Time, coordinates []float64, deviceID string) error {
	_, err := s.db.Exec("INSERT INTO traces (time, device_id, lon, lat) VALUES(?, ?, ?, ?)", t.Format(time.RFC3339), deviceID, coordinates[0], coordinates[1])
	return err
}

func (s *sqliteStore) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := s.db.QueryRowContext(ctx, "SELECT count(*) FROM traces").Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (s *sqliteStore) Open(ctx context.Context) error {
	db, err := sql.Open("sqlite3", s.path)
	if err != nil {
		return err
	}
	s.db = db
	return mig.Apply(ctx, db)
}

func (s *sqliteStore) Close(ctx context.Context) error {
	if s.db != nil {
		return nil
	}
	return s.db.Close()
}

func (s *sqliteStore) GetTracesSince(ctx context.Context, ts time.Time) ([]Trace, error) {
	start := ts.UTC().Format(time.RFC3339)
	rows, err := s.db.QueryContext(ctx, "SELECT time, lon, lat FROM traces where time >= ? order by time desc", start)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]Trace, 0, 100)
	for rows.Next() {
		tr := Trace{}
		if err := rows.Scan(&tr.Time, &tr.Lon, &tr.Lat); err != nil {
			return nil, err
		}
		result = append(result, tr)
	}
	return result, nil
}

func init() {
	mig = sqlitemigrate.NewRegistry()
	mig.RegisterMigration([]string{
		"CREATE TABLE IF NOT EXISTS traces (time text not null, device_id text, lon float, lat float)",
	}, []string{})
}
