package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Flags
	Json bool

	// Root Command
	rootCmd = &cobra.Command{
		Use:   "gitwho",
		Short: "Quickly get information about a github user.",
		Long:  "A CLI tool that allows you to quickly get information about a GitHub user, or search for someone even if you only know there name!",
		Run: func(cmd *cobra.Command, args []string) {
			if err := cmd.Usage(); err != nil {
				log.Fatal(err)
			}
		},
	}
)

func Execute() {
	// Global Flags
	rootCmd.PersistentFlags().BoolVarP(&Json, "json", "j", false, "Use this flag to dump the response to JSON output and a JSON file.")

	// Execute a command
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
