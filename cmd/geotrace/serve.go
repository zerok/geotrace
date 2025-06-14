package main

import (
	"fmt"
	"net/http"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/zerok/geotrace/internal/logging"
	"github.com/zerok/geotrace/pkg/server"
	"github.com/zerok/geotrace/pkg/store"
)

func generateServeCmd() *Command {
	var csvFile string
	var sqliteFile string
	var apiKey string
	var addr string
	var exposeMetrics bool
	cmd := &cobra.Command{
		Use: "serve",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			logger := logging.Setup()

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
			srv := server.New(st, apiKey, server.ExposeMetrics(exposeMetrics))
			if err := st.Open(ctx); err != nil {
				return err
			}
			defer st.Close(ctx)
			s.Handler = srv
			s.Addr = addr
			logger.InfoContext(ctx, "starting server", logging.Addr(s.Addr))
			return s.ListenAndServe()
		},
	}
	cmd.Flags().StringVar(&csvFile, "csv-store", "", "Path to a CSV file used for storage")
	cmd.Flags().StringVar(&sqliteFile, "sqlite-store", "", "Path to a SQLite file used for storage")
	cmd.Flags().StringVar(&apiKey, "api-key", "", "API key required to add new traces")
	cmd.Flags().StringVar(&addr, "addr", "localhost:8080", "Address to listen on")
	cmd.Flags().BoolVar(&exposeMetrics, "expose-metrics", false, "Expose metrics on /metrics")
	return &Command{cmd}
}
