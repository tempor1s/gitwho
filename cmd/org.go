package cmd

import (
	"context"
	"fmt"

	"github.com/google/go-github/v30/github"
	"github.com/logrusorgru/aurora"
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

	orgInfo := getOrgByName(orgName)

	if orgInfo == nil {
		fmt.Println("Error: Could not find Github Organization.")
		return
	}

	orgStruct := generateOrgStruct(orgInfo)

	printOrgInfo(orgStruct)
}

type githubOrg struct {
	Name        string
	Username    string
	Description string
	Location    string
	Website     string
	GithubUrl   string
	Email       string
	Repos       int
	Gists       int
	Followers   int
	Following   int
}

func generateOrgStruct(o *github.Organization) githubOrg {
	org := githubOrg{}

	org.Name = o.GetName()
	org.Username = o.GetLogin()
	org.Description = o.GetDescription()
	org.Location = o.GetLocation()
	org.Website = o.GetBlog()
	org.GithubUrl = o.GetHTMLURL()
	org.Email = o.GetEmail()
	org.Repos = o.GetPublicRepos()
	org.Gists = o.GetPublicGists()
	org.Followers = o.GetFollowers()
	org.Following = o.GetFollowing()

	return org
}

func printOrgInfo(o githubOrg) {
	fmt.Println(aurora.Bold("GitWho -- Simple GitHub information."))
	fmt.Println(aurora.Underline(aurora.Bold("General Info")))
	fmt.Println(aurora.Bold(aurora.Blue("- Org Name: ")), aurora.Bold(o.Name))
	fmt.Println(aurora.Bold(aurora.Blue("- Username: ")), aurora.Bold(o.Username))
	fmt.Println(aurora.Bold(aurora.Blue("- Description: ")), aurora.Bold(o.Description))
	fmt.Println(aurora.Bold(aurora.Blue("- Location: ")), aurora.Bold(o.Location))
	fmt.Println(aurora.Bold(aurora.Blue("- Website: ")), aurora.Bold(aurora.Underline(o.Website)))
	fmt.Println(aurora.Bold(aurora.Blue("- Link: ")), aurora.Bold(aurora.Underline(o.GithubUrl)))

	fmt.Println(aurora.Underline(aurora.Bold("By the numbers")))
	fmt.Println(aurora.Bold(aurora.Cyan("- Public Repos: ")), aurora.Bold(o.Repos))
	fmt.Println(aurora.Bold(aurora.Cyan("- Public Gists: ")), aurora.Bold(o.Gists))

	fmt.Println(aurora.Underline(aurora.Bold("Community")))
	fmt.Println(aurora.Bold(aurora.Green("- Followers: ")), aurora.Bold(o.Followers))
	fmt.Println(aurora.Bold(aurora.Green("- Following: ")), aurora.Bold(o.Following))
}

func getOrgByName(orgName string) *github.Organization {
	client := github.NewClient(nil)

	githubOrg, _, err := client.Organizations.Get(context.Background(), orgName)

	if err != nil {
		return nil
	}

	return githubOrg
}
