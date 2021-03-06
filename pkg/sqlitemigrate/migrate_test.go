package sqlitemigrate_test

import (
	"context"
	"database/sql"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
	"github.com/zerok/geotrace/pkg/sqlitemigrate"
)

func TestGetVersion(t *testing.T) {
	require.NoError(t, os.RemoveAll("data.sqlite"))
	ctx := context.Background()
	db, err := sql.Open("sqlite3", "./data.sqlite")
	require.NoError(t, err)
	_, err = db.Exec("PRAGMA user_version = 2")
	require.NoError(t, err)
	r := sqlitemigrate.NewRegistry()
	actual, err := r.GetCurrentVersion(ctx, db)
	require.NoError(t, err)
	require.Equal(t, 2, actual)
}

func TestMigrate(t *testing.T) {
	require.NoError(t, os.RemoveAll("data.sqlite"))
	ctx := context.Background()
	db, err := sql.Open("sqlite3", "./data.sqlite")
	require.NoError(t, err)
	_, err = db.Exec("PRAGMA user_version = 0")
	require.NoError(t, err)
	r := sqlitemigrate.NewRegistry()
	requireUserVersion(t, db, 0)

	r.RegisterMigration([]string{"CREATE TABLE table1 (id integer primary key autoincrement)"}, []string{"DROP TABLE table1"})
	require.NoError(t, r.Apply(ctx, db))
	requireUserVersion(t, db, 1)
}

func requireUserVersion(t *testing.T, db *sql.DB, expected int) {
	var actual int
	require.NoError(t, db.QueryRow("PRAGMA user_version").Scan(&actual))
	require.Equal(t, expected, actual)
}
