package kafka

import (
	"context"
	"log"
	"time"

	"github.com/harmonify/movie-reservation-system/cli/shared"

	"github.com/spf13/cobra"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type KafkaMigrateUpCmd struct {
	cmd        *cobra.Command
	path       string
	migrations []shared.KafkaMigration
	client     *kafka.AdminClient
	storage    shared.MigrationStorage
	logger     *log.Logger
}

func NewKafkaMigrateUpCmd(migrations []shared.KafkaMigration, client *kafka.AdminClient, storage shared.MigrationStorage, logger *log.Logger) *KafkaMigrateUpCmd {
	c := &KafkaMigrateUpCmd{
		path:       "root kafka migrate:up",
		migrations: migrations,
		client:     client,
		storage:    storage,
		logger:     logger,
	}
	c.cmd = &cobra.Command{
		Use:   "migrate:up",
		Short: "Apply pending Kafka migrations",
		Long:  "A CLI tool to manage Kafka migrations using Cobra, Viper, and SQLite.",
		Run: func(cmd *cobra.Command, args []string) {
			c.runMigrations()
		},
	}
	return c
}

func (c *KafkaMigrateUpCmd) Command() *cobra.Command {
	return c.cmd
}

func (c *KafkaMigrateUpCmd) Path() string {
	return c.path
}

func (c *KafkaMigrateUpCmd) runMigrations() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	state, err := c.storage.LoadState()
	if err != nil {
		log.Fatalf("Failed to load migration state: %v", err)
	}

	for _, migration := range c.migrations {
		identifier := migration.GetIdentifier()

		if applied, exists := state[identifier]; exists && applied {
			c.logger.Printf("Migration '%s' already applied, skipping", identifier)
			continue
		}

		c.logger.Printf("Applying migration '%s'", identifier)
		if err := migration.Up(ctx); err != nil {
			c.logger.Fatalf("Failed to apply migration '%s': %v", identifier, err)
		}

		if err := c.storage.SaveState(identifier, true); err != nil {
			c.logger.Fatalf("Failed to save migration state for '%s': %v", identifier, err)
		}
		c.logger.Printf("Migration '%s' applied successfully", identifier)
	}
}
