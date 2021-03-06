package store

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
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
func (s *csvFileStore) Count(ctx context.Context) (int64, error) {
	fp, err := s.fs.OpenFile(s.path, os.O_CREATE|os.O_RDONLY, 0600)
	if err != nil {
		return 0, err
	}
	defer fp.Close()
	reader := csv.NewReader(fp)
	reader.Comma = ';'
	var count int64
	for {
		_, err := reader.Read()
		if io.EOF == err {
			break
		}
		if err != nil {
			return 0, err
		}
		count += 1
	}
	return count, nil
}

func (s *csvFileStore) Close(ctx context.Context) error {
	return nil
}

func (s *csvFileStore) Open(ctx context.Context) error {
	return nil
}

func (s *csvFileStore) GetTracesSince(ctx context.Context, ts time.Time) ([]Trace, error) {
	return nil, fmt.Errorf("not supported")
}
