package kafka_migration

const MigrationTemplate = `package kafka_migration

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type {{.PascalCasedName}}Migration struct {
	client *kafka.AdminClient
}

func New{{.PascalCasedName}}Migration(client *kafka.AdminClient) *{{.PascalCasedName}}Migration {
	return &{{.PascalCasedName}}Migration{
		client: client,
	}
}

func (m *{{.PascalCasedName}}Migration) GetIdentifier() string {
	return "{{.Timestamp}}_{{.Name}}"
}

func (m *{{.PascalCasedName}}Migration) Up(ctx context.Context) error {
	panic("{{.PascalCasedName}}Migration.Up not implemented")
}

func (m *{{.PascalCasedName}}Migration) Down(ctx context.Context) error {
	panic("{{.PascalCasedName}}Migration.Down not implemented")
}
`

type MigrationTemplateData struct {
	Name            string
	PascalCasedName string
	Timestamp       string
}
