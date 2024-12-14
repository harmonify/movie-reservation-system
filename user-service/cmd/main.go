package cmd

import (
	"log"
	"os"
	"github.com/spf13/cobra"
)

var (
	env       string
	dbAction  string
	syncForce bool
)

var rootCmd = &cobra.Command{
	Use:   "harmonify-utility",
	Short: "Command Utility Tools for harmonify Project",
	Long:  `Command Utility Tools for harmonify Projects. You can execute any command related to this project. For example: Start Service with specific environment, or execute data seeder, run the test.etc`,
}

var startCmd = &cobra.Command{
	Use:       "start",
	Short:     "Start the server with specific environment: [ dev | prod ]",
	Long:      `Start the server with specific environment: [ dev | prod ]. By default it will use dev environment`,
	ValidArgs: []string{"dev", "prod"},
	Run: func(cmd *cobra.Command, args []string) {
		for _, v := range cmd.ValidArgs {
			if v == env {
				err := os.Setenv("ENV", env)
				if err != nil {
					log.Panic(err)
				}
				newServer()
			}
		}
		log.Panic("Env value only valid with: ", cmd.ValidArgs, ". given value:", env)
	},
}

var testCmd = &cobra.Command{
	Use:       "test",
	Short:     "Run the test",
	Long:      `Run the test.`,
	ValidArgs: []string{"dev", "prod"},
	Run: func(cmd *cobra.Command, args []string) {
		for _, v := range cmd.ValidArgs {
			if v == env {
				err := os.Setenv("ENV", env)
				if err != nil {
					log.Panic(err)
				}
				pkg.Server()
			}
		}
		log.Panic("Env value only valid with: ", cmd.ValidArgs, ". given value:", env)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.Flags().StringVarP(&env, "env", "e", "dev", "Environment: dev | prod")
}
