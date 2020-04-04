package cmd

import (
	"context"
	"fmt"

	"github.com/google/go-github/v30/github"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(orgCommand)
}

var (
	// Flags
	// TODO: User flag to get all user info from an org ;)

	// Cmd
	orgCommand = &cobra.Command{
		Use:   "org",
		Short: "Get information about a specific Github org.",
		Long:  "Allows you to get information about a specific GitHub organization, but just its name.. Example: `gitwho org google`",
		Run:   orgCmd,
	}
)

func orgCmd(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		fmt.Println("Please enter the org you want to get information about. `gitwho org google`")
		return
	}

	orgName := args[0]

	fmt.Println("Org Name:", orgName)

	githubOrg := getOrgByName(orgName)

	if githubOrg == nil {
		fmt.Println("Error: Could not find Github Organization.")
		return
	}

	fmt.Printf("%+v\n", githubOrg)
}

func getOrgByName(orgName string) *github.Organization {
	client := github.NewClient(nil)

	githubOrg, _, err := client.Organizations.Get(context.Background(), orgName)

	if err != nil {
		return nil
	}

	return githubOrg
}
