package shared

import (
	"context"

	"github.com/spf13/cobra"
)

type KafkaMigration interface {
	// GetIdentifier returns unique identifier for the migration
	GetIdentifier() string
	// Up operation must be idempotent
	Up(ctx context.Context) error
	// Down operation must be idempotent
	Down(ctx context.Context) error
}

type MigrationStorage interface {
	LoadState() (map[string]bool, error)
	SaveState(identifier string, applied bool) error
}

type CobraCommand interface {
	Command() *cobra.Command
	Path() string
}
