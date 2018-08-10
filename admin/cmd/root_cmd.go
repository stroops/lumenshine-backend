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
	return &rootCmd
}

func run(cmd *cobra.Command, args []string) {
	//fmt.Println("args[0]: " + args[0])
}
