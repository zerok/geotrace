package store

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"

	"github.com/spf13/afero"
)

type csvFileStore struct {
	fs   afero.Fs
	path string
}

func NewCSVFileStore(fs afero.Fs, path string) Store {
	return &csvFileStore{fs: fs, path: path}
}

func (s *csvFileStore) Add(t time.Time, coordinates []float64, deviceID string) error {
	fp, err := s.fs.OpenFile(s.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer fp.Close()
	writer := csv.NewWriter(fp)
	writer.Comma = ';'
	defer writer.Flush()
	if err := writer.Write([]string{
		t.Format(time.RFC3339),
		fmt.Sprintf("%f %f", coordinates[0], coordinates[1]),
		deviceID,
	}); err != nil {
		return err
	}
	return nil
}
