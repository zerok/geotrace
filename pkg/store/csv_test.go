package store_test

import (
	"testing"
	"time"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
	"github.com/zerok/geotrace/pkg/store"
)

func TestCSVFileStore(t *testing.T) {
	fs := &afero.MemMapFs{}
	path := "dummy.csv"
	now := time.Date(2020, time.May, 10, 12, 13, 14, 0, time.UTC)
	s := store.NewCSVFileStore(fs, path)
	require.NoError(t, s.Add(now, []float64{1, 2}, ""))
	data, err := afero.ReadFile(fs, path)
	require.NoError(t, err)
	require.Equal(t, "2020-05-10T12:13:14Z;1.000000 2.000000;\n", string(data))
}
