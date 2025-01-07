package main

import (
	"context"
	"kafka-playground/cmd"
	"kafka-playground/migrations/kafka"
	"kafka-playground/shared"
	"log"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

func NewApp(p ...fx.Option) *fx.App {
	options := []fx.Option{
		fx.Provide(
			func() *log.Logger {
				return log.Default()
			},
			func() *shared.ConfigFile {
				return &shared.ConfigFile{
					Path: ".env",
				}
			},
		),
		kafka.MigrationModule,
		shared.SharedModule,
		cmd.CommandModule,

		fx.NopLogger, // Removes all fx logs, even on error
		fx.Invoke(func(*cobra.Command) {}),
	}

	if len(p) > 0 {
		options = append(options, p...)
	}

	return fx.New(options...)
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	if err := NewApp().Start(ctx); err != nil {
		log.Fatalf("Failed to start app: %v", err)
	}
}
