package main

import (
	"fmt"
	"net/http"

	"github.com/rs/zerolog"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/zerok/geotrace/pkg/server"
	"github.com/zerok/geotrace/pkg/store"
)

func generateServeCmd() *Command {
	var csvFile string
	var sqliteFile string
	var apiKey string
	var addr string
	cmd := &cobra.Command{
		Use: "serve",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			logger := zerolog.Ctx(ctx)

			if csvFile == "" && sqliteFile == "" {
				return fmt.Errorf("specify a CSV or SQLite file")
			}
			s := http.Server{}
			var st store.Store
			if sqliteFile != "" {
				st = store.NewSQLiteStore(sqliteFile)
			} else {
				st = store.NewCSVFileStore(afero.NewOsFs(), csvFile)
			}
			srv := server.New(st, apiKey)
			if err := st.Open(ctx); err != nil {
				return err
			}
			defer st.Close(ctx)
			s.Handler = srv
			s.Addr = addr
			logger.Info().Msgf("Starting server on %s", s.Addr)
			return s.ListenAndServe()
		},
	}
	cmd.Flags().StringVar(&csvFile, "csv-store", "", "Path to a CSV file used for storage")
	cmd.Flags().StringVar(&sqliteFile, "sqlite-store", "", "Path to a SQLite file used for storage")
	cmd.Flags().StringVar(&apiKey, "api-key", "", "API key required to add new traces")
	cmd.Flags().StringVar(&addr, "addr", "localhost:8080", "Address to listen on")
	return &Command{cmd}
}
