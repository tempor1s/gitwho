package cmd

import (
	"context"
	"fmt"
	"time"

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

	// Gitub user
	userinfo := generateUserStruct(githubUser)

	// Print info.
	printUserInfo(userinfo)
}

type githubUser struct {
	Name           string
	Username       string
	Bio            string
	Location       string
	Website        string
	GithubUrl      string
	Hireable       bool
	Org            string
	Repos          int
	Gists          int
	Followers      int
	Following      int
	LastActive     time.Time
	AccountCreated time.Time
}

// generateuserStruct cleans the response from github and turns it into a nice struct for us to use
func generateUserStruct(gu *github.User) githubUser {
	user := githubUser{}

	user.Name = gu.GetName()
	user.Username = gu.GetLogin()
	user.Bio = gu.GetBio()
	user.Location = gu.GetLocation()
	user.Website = gu.GetBlog()
	user.GithubUrl = gu.GetHTMLURL()
	user.Hireable = gu.GetHireable()
	user.Org = gu.GetCompany()
	user.Repos = gu.GetPublicRepos()
	user.Gists = gu.GetPublicGists()
	user.Followers = gu.GetFollowers()
	user.Following = gu.GetFollowing()
	user.LastActive = gu.GetUpdatedAt().Time
	user.AccountCreated = gu.GetCreatedAt().Time

	return user
}

// printUserInfo will print out our nicely formatted struct :)
func printUserInfo(u githubUser) {
	fmt.Println("General Info")
	fmt.Printf("- Real Name: %s\n", u.Name)
	fmt.Printf("- Username: %s\n", u.Username)
	fmt.Printf("- Bio: %s\n", u.Bio)
	fmt.Printf("- Location: %s\n", u.Location)
	fmt.Printf("- Website: %s\n", u.Website)
	fmt.Printf("- Link: %s\n", u.GithubUrl)

	fmt.Println("Work Info")
	fmt.Printf("- Hireable: %t\n", u.Hireable)
	fmt.Printf("- Organization: %s\n", u.Org)
	fmt.Println("By the numbers")
	fmt.Printf("- Public Repos: %d\n", u.Repos)
	fmt.Printf("- Public Gists: %d\n", u.Gists)
	fmt.Println("Community")
	fmt.Printf("- Followers: %d\n", u.Followers)
	fmt.Printf("- Following: %d\n", u.Following)

	fmt.Println("Dates")
	year, month, day := u.LastActive.Date()
	fmt.Printf("- Last Active: %s (%d/%d/%d)\n", humanize.Time(u.LastActive), month, day, year)
	year, month, day = u.AccountCreated.Date()
	fmt.Printf("- Account Created: %s (%d/%d/%d)\n", humanize.Time(u.AccountCreated), month, day, year)
}

// getUserByUsername gets user info from the github API
func getUserByUsername(username string) *github.User {
	client := github.NewClient(nil)

	githubUser, _, err := client.Users.Get(context.Background(), username)

	if err != nil {
		return nil
	}

	return githubUser
}
