package cmd

import (
	"context"
	"log"
	"os"

	"github.com/google/go-github/v30/github"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

var (
	// Flags
	Json  bool
	Token string

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
	rootCmd.PersistentFlags().StringVarP(&Token, "token", "t", "", "Supply a Personal Access Tokent to up the api request limit and access private information.")

	// Execute a command
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func createGithubClient() *github.Client {
	// Set up OAuth token stuff
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: Token},
	)

	tc := oauth2.NewClient(ctx, ts)

	// Create a new github client using the OAuth2 token or no token
	var client *github.Client
	if Token != "" {
		client = github.NewClient(tc)
	} else {
		client = github.NewClient(nil)
	}

	return client
}
