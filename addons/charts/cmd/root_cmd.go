package cmd

// this is cmd/root_cmd.go

import (
	"github.com/spf13/cobra"
)

// RootCommand will setup and return the root command
func RootCommand() *cobra.Command {
	rootCmd := cobra.Command{
		Use: "dbCommand",
		Run: run,
	}

	// this is where we will configure everything!
	rootCmd.Flags().StringP("config-dir", "c", "", "data directory with config.toml e.g. /opt/wwwdata/if6/mkt/api/data")
	rootCmd.Flags().Bool("history", false, "set to true in order to generate the history table")
	rootCmd.Flags().Bool("truncate", false, "set true to truncate history table. only works if history flag is true")

	return &rootCmd
}

func run(cmd *cobra.Command, args []string) {
	// fmt.Println("--- here ---")
}
