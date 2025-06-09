package main

import (
	"context"
	"os"

	"github.com/spf13/cobra"
	"github.com/zerok/geotrace/internal/logging"
)

type Command struct {
	*cobra.Command
}

func generateRootCmd() *Command {
	cmd := cobra.Command{
		Use: "geotrace",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	cmd.AddCommand(generateServeCmd().Command)
	cmd.AddCommand(generateExportTrackCmd().Command)
	cmd.AddCommand(generateBrowseCmd().Command)
	return &Command{&cmd}
}

func main() {
	ctx := context.Background()
	cmd := generateRootCmd()
	logger := logging.Setup()
	if err := cmd.ExecuteContext(ctx); err != nil {
		logger.ErrorContext(ctx, "command failed", logging.Err(err))
		os.Exit(1)
	}
}
