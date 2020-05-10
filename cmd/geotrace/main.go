package main

import (
	"context"
	"os"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
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
	return &Command{&cmd}
}

func main() {
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout})
	cmd := generateRootCmd()
	if err := cmd.ExecuteContext(logger.WithContext(context.Background())); err != nil {
		logger.Fatal().Err(err).Msg("Command failed.")
	}
}
