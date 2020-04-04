package cmd

import (
	"context"
	"fmt"

	"github.com/dustin/go-humanize"
	"github.com/google/go-github/v30/github"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(usernameCommand)
}

var usernameCommand = &cobra.Command{
	Use:   "username",
	Short: "Get information about a specific Github user. ",
	Long:  "Allows you to get information about a specific GitHub user just from their github username. Example: `gitwho username tempor1s`",
	Run:   usernameCmd,
}

func usernameCmd(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		fmt.Println("Please pass in the username you want to search up. Example: `gitwho username tempor1s`")
		return
	}

	username := args[0]

	githubUser := getUserByUsername(username)

	if githubUser == nil {
		fmt.Println("Error: User not found.")
		return
	}

	// fmt.Printf("- User: %v\n", githubUser)
	fmt.Println("General Info")
	fmt.Printf("- Real Name: %s\n", githubUser.GetName())
	fmt.Printf("- Username: %s\n", githubUser.GetLogin())
	fmt.Printf("- Location: %s\n", githubUser.GetLocation())

	fmt.Println("Work Info")
	fmt.Printf("- Hireable: %t\n", githubUser.GetHireable())
	fmt.Printf("- Website: %s\n", githubUser.GetBlog())
	fmt.Println("By the numbers")
	fmt.Printf("- Public Repos: %d\n", githubUser.GetPublicRepos())
	fmt.Printf("- Public Gists: %d\n", githubUser.GetPublicGists())
	fmt.Println("Community")
	fmt.Printf("- Followers: %d\n", githubUser.GetFollowers())
	fmt.Printf("- Following: %d\n", githubUser.GetFollowing())

	fmt.Println("Dates")
	fmt.Printf("- Last Active: %s\n", humanize.Time(githubUser.UpdatedAt.Time))
	year, month, day := githubUser.CreatedAt.Date()
	fmt.Printf("- Account Created: %d/%d/%d\n", month, day, year)
}

func getUserByUsername(username string) *github.User {
	client := github.NewClient(nil)

	githubUser, _, err := client.Users.Get(context.Background(), username)

	if err != nil {
		return nil
	}

	return githubUser
}
