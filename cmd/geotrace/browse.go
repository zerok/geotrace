package main

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/zerok/geotrace/internal/logging"
	"github.com/zerok/geotrace/pkg/browser"
	"github.com/zerok/geotrace/pkg/store"
)

func generateBrowseCmd() *Command {
	var sqliteFile string
	var addr string
	var webRoot string
	cmd := &cobra.Command{
		Use: "browse",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			logger := logging.Setup()

			if sqliteFile == "" {
				return fmt.Errorf("specify a SQLite file")
			}
			s := http.Server{}
			st := store.NewSQLiteStore(sqliteFile)
			srv := browser.New(st, browser.WithWebRoot(webRoot))
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
	cmd.Flags().StringVar(&sqliteFile, "sqlite-store", "", "Path to a SQLite file used for storage")
	cmd.Flags().StringVar(&addr, "addr", "localhost:8080", "Address to listen on")
	cmd.Flags().StringVar(&webRoot, "webroot", "", "Path to a folder that should be exposed via /")
	return &Command{cmd}
}
