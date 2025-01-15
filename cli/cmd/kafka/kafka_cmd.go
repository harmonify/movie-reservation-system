package kafka

import "github.com/spf13/cobra"

type KafkaCmd struct {
	cmd  *cobra.Command
	path string
}

func (c KafkaCmd) Command() *cobra.Command {
	return c.cmd
}

func (c KafkaCmd) Path() string {
	return c.path
}

func NewKafkaCmd() *KafkaCmd {
	return &KafkaCmd{
		cmd: &cobra.Command{
			Use:   "kafka",
			Short: "Kafka utility CLI",
			Run: func(cmd *cobra.Command, args []string) {
				cmd.Help()
			},
		},
		path: "root kafka",
	}
}
