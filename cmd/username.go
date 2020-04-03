package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(usernameCommand)
}

var usernameCommand = &cobra.Command{
	Use:   "username",
	Short: "Get information about a specific Github user. ",
	Long:  "Allows you to get information about a specific GitHub user just from their github username. Example: `gitwho username tempor1s`",
	Run:   userLookup,
}

func userLookup(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		fmt.Println("Please pass in the username you want to search up. Example: `gitwho username tempor1s`")
		return
	}

	username := args[0]

	fmt.Printf("Username: %s", username)
}
