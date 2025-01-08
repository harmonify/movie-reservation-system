package kafka

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/harmonify/movie-reservation-system/cli/shared"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/spf13/cobra"
)

type KafkaMigrateDownCmd struct {
	cmd        *cobra.Command
	path       string
	migrations []shared.KafkaMigration
	client     *kafka.AdminClient
	storage    shared.MigrationStorage
	logger     *log.Logger
}

func NewKafkaMigrateDownCmd(migrations []shared.KafkaMigration, client *kafka.AdminClient, storage shared.MigrationStorage, logger *log.Logger) *KafkaMigrateDownCmd {
	c := &KafkaMigrateDownCmd{
		path:       "root kafka migrate:down",
		migrations: migrations,
		client:     client,
		storage:    storage,
		logger:     logger,
	}
	c.cmd = &cobra.Command{
		Use:   "migrate:down [steps]",
		Short: "Revert Kafka migrations",
		Long:  "A CLI tool to revert Kafka migrations by a specified number of steps.",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var steps int
			if len(args) > 0 {
				var err error
				steps, err = strconv.Atoi(args[0])
				if err != nil {
					log.Fatalf("Invalid steps argument: %v", err)
				}
			}
			c.revertMigrations(steps)
		},
	}
	return c
}

func (c *KafkaMigrateDownCmd) Command() *cobra.Command {
	return c.cmd
}

func (c *KafkaMigrateDownCmd) Path() string {
	return c.path
}

func (c *KafkaMigrateDownCmd) revertMigrations(steps int) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	state, err := c.storage.LoadState()
	if err != nil {
		log.Fatalf("Failed to load migration state: %v", err)
	}

	appliedMigrations := make([]shared.KafkaMigration, 0)
	for _, migration := range c.migrations {
		identifier := migration.GetIdentifier()
		if applied, exists := state[identifier]; exists && applied {
			appliedMigrations = append(appliedMigrations, migration)
		}
	}

	// Revert specified number of migrations, or all if steps is 0
	toRevert := appliedMigrations
	if steps > 0 && steps < len(appliedMigrations) {
		toRevert = appliedMigrations[len(appliedMigrations)-steps:]
	}

	for i := len(toRevert) - 1; i >= 0; i-- { // Reverse order for downward migrations
		migration := toRevert[i]
		identifier := migration.GetIdentifier()

		c.logger.Printf("Reverting migration '%s'", identifier)
		if err := migration.Down(ctx); err != nil {
			c.logger.Fatalf("Failed to revert migration '%s': %v", identifier, err)
		}

		if err := c.storage.SaveState(identifier, false); err != nil {
			c.logger.Fatalf("Failed to update migration state for '%s': %v", identifier, err)
		}
		c.logger.Printf("Migration '%s' reverted successfully", identifier)
	}
}
