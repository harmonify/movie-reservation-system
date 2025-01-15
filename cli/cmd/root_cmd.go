package cmd

import (
	"context"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/harmonify/movie-reservation-system/cli/shared"

	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

func NewRootCmd(commands []shared.CobraCommand, lc fx.Lifecycle, shutdowner fx.Shutdowner) *cobra.Command {
	root := &cobra.Command{
		Use:   "mrs-cli",
		Short: "Movie reservation system CLI",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	// Loop through commands adding them to other commands
	for _, cmd := range commands {
		parentalPath := cmd.Path()[:strings.LastIndex(cmd.Path(), " ")]
		if parentalPath == "root" {
			root.AddCommand(cmd.Command())
		} else {
			c2 := commands[slices.IndexFunc(commands, func(c shared.CobraCommand) bool { return c.Path() == parentalPath })]
			c2.Command().AddCommand(cmd.Command())
		}
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if err := root.Execute(); err != nil {
				fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing the CLI '%s'", err)
				shutdowner.Shutdown(fx.ExitCode(1))
			}
			shutdowner.Shutdown()
			return nil
		},
	})

	return root
}
