# CLI

The Movie Reservation System includes a robust Command Line Interface (CLI) that facilitates seamless interaction with the system. It is designed for developers and system administrators to manage critical operations such as Kafka topic migrations and configuration.

## Kafka Migrations

Kafka migrations are a critical component of the Movie Reservation System, allowing smooth and controlled creation, update, and deletion of Kafka topics.

The CLI leverages industry-standard tools such as **Cobra** for CLI framework, **Viper** for configuration management, **FX** for dependency injection, and **SQLite** for persistent storage.

The CLI provides the following commands under the `kafka` namespace:

### `kafka migrate:up`

Applies all pending Kafka topic migrations. This command ensures that the Kafka infrastructure is kept up-to-date and aligned with application requirements. It supports:

-   **Idempotency:** Each migration is applied only once, ensuring safe re-execution of the command.
-   **Persistence:** Migration states are stored locally in SQLite and optionally persisted to S3-compatible storage for disaster recovery and distributed consistency.
-   **Timeout Management:** Configurable execution timeouts prevent long-running or hanging processes.

Usage:

```bash
mrs-cli kafka migrate:up
```

Steps Executed:

-   Loads the current migration state from SQLite.
-   Executes all unapplied migrations using the Up method.
-   Updates the SQLite database.
-   Persists state to S3 storage upon successful completion. (TODO)

### `kafka migrate:down`

Rolls back Kafka migrations. The command can optionally accept a number to revert a specified count of recent migrations. It supports:

-   Selective Rollback: Supports partial migration rollbacks for testing and debugging.
-   Idempotency: Ensures safe and consistent rollback of migrations.

Usage:

```bash
mrs-cli kafka migrate:down [count]
```

Steps Executed:

-   Loads the current migration state from SQLite.
-   Rolls back migrations using the Down method, starting from the latest applied migration.
-   Updates the SQLite database.
-   Persists state to S3 storage upon successful completion. (TODO)

### Example Workflow

#### Step 1: Create a New Migration

Add a new migration by implementing the KafkaMigration interface. Example:

```go
type CreateTopicMigration struct {
    // ...
}

func (m *CreateTopicMigration) GetIdentifier() string {
    return "create-new-topic"
}

func (m *CreateTopicMigration) Up(ctx context.Context) error {
    // Logic to create Kafka topic
}

func (m *CreateTopicMigration) Down(ctx context.Context) error {
    // Logic to delete Kafka topic
}
```

#### Step 2: Run Migrations

Apply all pending migrations:

```bash
mrs-cli kafka migrate:up
```

#### Step 3: Rollback Migrations

Rollback the last migration:

```bash
mrs-cli kafka migrate:down 1
```
