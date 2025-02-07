package kafka

import (
	"log"
	"os"
	"path"
	"runtime"
	"strings"
	"text/template"
	"time"

	kafka_migration "github.com/harmonify/movie-reservation-system/cli/migrations/kafka"
	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type KafkaMigrateNewCmd struct {
	cmd    *cobra.Command
	path   string
	logger *log.Logger
}

func NewKafkaMigrateNewCmd(logger *log.Logger) *KafkaMigrateNewCmd {
	c := &KafkaMigrateNewCmd{
		path:   "root kafka migrate:new",
		logger: logger,
	}
	c.cmd = &cobra.Command{
		Use:   "migrate:new [name]",
		Short: "Create new Kafka migration",
		Long:  "A CLI tool to create new Kafka migration with specified name.",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				c.logger.Fatalf("Migration name is required")
				c.logger.Printf("Example usage: %s", "migrate:new create_public.user.registered.v1_topic")
			}
			c.createNewMigration(args[0])
		},
	}
	return c
}

func (c *KafkaMigrateNewCmd) Command() *cobra.Command {
	return c.cmd
}

func (c *KafkaMigrateNewCmd) Path() string {
	return c.path
}

func (c *KafkaMigrateNewCmd) createNewMigration(name string) {
	now := getTimestamp()

	migration, err := c.renderNewMigration(name, now)
	if err != nil {
		c.logger.Fatalf("Error rendering new migration: %v", err)
	}

	// get current file path
	_, file, _, _ := runtime.Caller(0)

	// create new migration file
	newMigrationPath := path.Join(path.Dir(file), "..", "..", "migrations", "kafka", now+"_"+name+".go")
	f, err := os.Create(newMigrationPath)
	if err != nil {
		c.logger.Fatalf("Error creating new migration file: %v", err)
		return
	}
	defer f.Close()

	// write migration to file
	_, err = f.WriteString(migration)
	if err != nil {
		c.logger.Fatalf("Error writing migration to file: %v", err)
		return
	}

	c.logger.Printf("Migration created: %s", newMigrationPath)
	c.logger.Printf("Register the migration in the %s", path.Join(path.Dir(file), "..", "..", "migrations", "kafka", "module.go"))
}

// renderNewMigration returns a new migration template
// with the given name.
// i.e. renderNewMigration("create_public.user.registered.v1_topic", "20210102150405")
func (c *KafkaMigrateNewCmd) renderNewMigration(name string, timestamp string) (string, error) {
	pascalCasedName := anyToPascal(name)

	tmpl, err := template.New("migration").Parse(kafka_migration.MigrationTemplate)
	if err != nil {
		return "", err
	}

	var buf strings.Builder
	err = tmpl.Execute(&buf, kafka_migration.MigrationTemplateData{
		Name:            name,
		PascalCasedName: pascalCasedName,
		Timestamp:       timestamp,
	})
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

// Split the name by ".", "_", "-", and convert each part to PascalCase
func anyToPascal(name string) string {
	pascal := ""
	parts := strings.FieldsFunc(name, func(r rune) bool {
		return r == '.' || r == '_' || r == '-'
	})
	for _, part := range parts {
		pascal += cases.Title(language.English).String(part)
	}
	return pascal
}

// getTimestamp returns the current timestamp in the following example format: 20060102150405
func getTimestamp() string {
	return time.Now().Format("20060102150405")
}
