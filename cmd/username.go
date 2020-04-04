package cmd

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/google/go-github/v30/github"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

func init() {
	// register the command
	rootCmd.AddCommand(usernameCommand)
	// setup local flags
	usernameCommand.Flags().BoolVarP(&Open, "open", "o", false, "Open their GitHub repo after printing info.")
}

var (
	// local flags
	Open bool
	// local cmd
	usernameCommand = &cobra.Command{
		Use:   "username",
		Short: "Get information about a specific Github user. ",
		Long:  "Allows you to get information about a specific GitHub user just from their github username. Example: `gitwho username tempor1s`",
		Run:   usernameCmd,
	}
)

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

	if Open {
		err := exec.Command("open", userinfo.GithubUrl).Start()

		if err != nil {
			fmt.Println(aurora.Bold(aurora.Red("Error: Could not open Github URL: ")), userinfo.GithubUrl)
			return
		}
	}
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
	fmt.Println(aurora.Bold("GitWho -- Simple GitHub information."))
	fmt.Println(aurora.Underline(aurora.Bold("General Info")))
	fmt.Println(aurora.Bold(aurora.Blue("- Real Name: ")), aurora.Bold(u.Name))
	fmt.Println(aurora.Bold(aurora.Blue("- Username: ")), aurora.Bold(u.Username))
	fmt.Println(aurora.Bold(aurora.Blue("- Bio: ")), aurora.Bold(u.Bio))
	fmt.Println(aurora.Bold(aurora.Blue("- Location: ")), aurora.Bold(u.Location))
	fmt.Println(aurora.Bold(aurora.Blue("- Website: ")), aurora.Bold(aurora.Underline(u.Website)))
	fmt.Println(aurora.Bold(aurora.Blue("- Link: ")), aurora.Bold(aurora.Underline(u.GithubUrl)))

	fmt.Println(aurora.Underline(aurora.Bold("Work Info")))
	fmt.Println(aurora.Bold(aurora.Magenta("- Hireable")), aurora.Bold(u.Hireable))
	fmt.Println(aurora.Bold(aurora.Magenta("- Organization: ")), aurora.Bold(u.Org))

	fmt.Println(aurora.Underline(aurora.Bold("By the numbers")))
	fmt.Println(aurora.Bold(aurora.Cyan("- Public Repos: ")), aurora.Bold(u.Repos))
	fmt.Println(aurora.Bold(aurora.Cyan("- Public Gists: ")), aurora.Bold(u.Gists))

	fmt.Println(aurora.Underline(aurora.Bold("Community")))
	fmt.Println(aurora.Bold(aurora.Green("- Followers: ")), aurora.Bold(u.Followers))
	fmt.Println(aurora.Bold(aurora.Green("- Following: ")), aurora.Bold(u.Following))

	fmt.Println(aurora.Underline(aurora.Bold("Dates")))
	year, month, day := u.LastActive.Date()
	fmt.Println(aurora.Bold(aurora.Yellow("- Last Active: ")), aurora.Bold(fmt.Sprintf("%s (%d/%d/%d)", humanize.Time(u.LastActive), month, day, year)))
	year, month, day = u.AccountCreated.Date()
	fmt.Println(aurora.Bold(aurora.Yellow("- Account Created: ")), aurora.Bold(fmt.Sprintf("%s (%d/%d/%d)", humanize.Time(u.AccountCreated), month, day, year)))
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
