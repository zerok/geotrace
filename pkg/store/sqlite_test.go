package store_test

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/stretchr/testify/require"
	"github.com/zerok/geotrace/pkg/store"
)

func TestSQLiteStore(t *testing.T) {
	ctx := context.Background()
	os.RemoveAll("data.sqlite")
	st := store.NewSQLiteStore("data.sqlite")
	require.NotNil(t, st)
	require.NoError(t, st.Open(ctx))
	defer st.Close(ctx)
	now := time.Date(2020, 5, 11, 12, 13, 14, 0, time.UTC)
	lat := 48.487486
	lon := 13.046261
	require.NoError(t, st.Add(now, []float64{lon, lat}, ""))

	db, err := sql.Open("sqlite3", "data.sqlite")
	require.NoError(t, err)
	defer db.Close()
	var count int
	var retrievedLat float64
	var retrievedLon float64
	require.NoError(t, db.QueryRow("SELECT count(*) FROM traces").Scan(&count))
	require.Equal(t, 1, count)

	require.NoError(t, db.QueryRow("SELECT lat, lon FROM traces").Scan(&retrievedLat, &retrievedLon))
	require.Equal(t, lat, retrievedLat)
	require.Equal(t, lon, retrievedLon)
}
