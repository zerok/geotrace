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
	var apiKey string
	var addr string
	cmd := &cobra.Command{
		Use: "serve",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			logger := zerolog.Ctx(ctx)

			if csvFile == "" {
				return fmt.Errorf("specify a CSV file")
			}
			s := http.Server{}
			st := store.NewCSVFileStore(afero.NewOsFs(), csvFile)
			srv := server.New(st, apiKey)
			s.Handler = srv
			s.Addr = addr
			logger.Info().Msgf("Starting server on %s", s.Addr)
			return s.ListenAndServe()
		},
	}
	cmd.Flags().StringVar(&csvFile, "csv-store", "", "Path to a CSV file used for storage")
	cmd.Flags().StringVar(&apiKey, "api-key", "", "API key required to add new traces")
	cmd.Flags().StringVar(&addr, "addr", "localhost:8080", "Address to listen on")
	return &Command{cmd}
}
