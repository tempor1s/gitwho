package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/google/go-github/v30/github"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(orgCommand)

	orgCommand.Flags().BoolVarP(&Users, "users", "u", false, "Use this flag if you would like to gather all users information for the org. Note: This might take awhile!")
}

var (
	// Flags
	Users bool

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
	orgMembers := getOrgMembers(orgName)

	if orgInfo == nil {
		fmt.Println("Error: Could not find Github Organization.")
		return
	}

	orgStruct := generateOrgStruct(orgInfo, orgMembers)

	if Json {
		printJsonOrgInfo(orgStruct)
	} else {
		printOrgInfo(orgStruct)
	}
}

type githubOrg struct {
	Name              string       `json:"name"`
	Username          string       `json:"org_name"`
	Description       string       `json:"description"`
	Location          string       `json:"location"`
	Website           string       `json:"website"`
	GithubUrl         string       `json:"github_url"`
	Email             string       `json:"email"`
	Repos             int          `json:"repos"`
	Gists             int          `json:"gists"`
	Followers         int          `json:"followers"`
	Following         int          `json:"following"`
	PublicMemberCount int          `json:"public_member_count"`
	PublicMembers     []githubUser `json:"org_members"`
}

func generateOrgStruct(o *github.Organization, orgMembers []*github.User) githubOrg {
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
	org.PublicMemberCount = len(orgMembers)

	publicMembers := []githubUser{}

	// Convert all *githubUser to our custom struct - this could take awhile..
	if Users {
		for _, member := range orgMembers {
			fmt.Printf("%+v\n", member)

			githubUser := getUserByUsername(member.GetLogin())
			usr := generateUserStruct(githubUser)
			publicMembers = append(publicMembers, usr)
		}

		org.PublicMembers = publicMembers
	}

	return org
}

func printJsonOrgInfo(o githubOrg) {
	json, err := json.MarshalIndent(&o, "", "    ")

	if err != nil {
		log.Fatal(err)
	}

	ioutil.WriteFile(fmt.Sprintf("%s.json", o.Username), json, os.ModePerm)

	fmt.Printf("%s\n", string(json))
}

func printOrgInfo(o githubOrg) {
	fmt.Println(aurora.Bold("GitWho -- Simple GitHub information."))
	fmt.Println(aurora.Underline(aurora.Bold("General Info")))
	fmt.Println(aurora.Bold(aurora.Blue("- Org Name:")), aurora.Bold(o.Name))
	fmt.Println(aurora.Bold(aurora.Blue("- Username:")), aurora.Bold(o.Username))
	fmt.Println(aurora.Bold(aurora.Blue("- Description:")), aurora.Bold(o.Description))
	fmt.Println(aurora.Bold(aurora.Blue("- Location:")), aurora.Bold(o.Location))
	fmt.Println(aurora.Bold(aurora.Blue("- Website:")), aurora.Bold(aurora.Underline(o.Website)))
	fmt.Println(aurora.Bold(aurora.Blue("- Link:")), aurora.Bold(aurora.Underline(o.GithubUrl)))

	fmt.Println(aurora.Underline(aurora.Bold("By the numbers")))
	fmt.Println(aurora.Bold(aurora.Cyan("- Public Repos:")), aurora.Bold(o.Repos))
	fmt.Println(aurora.Bold(aurora.Cyan("- Public Gists:")), aurora.Bold(o.Gists))

	fmt.Println(aurora.Underline(aurora.Bold("Community")))
	msg := ""
	if o.PublicMemberCount == 100 {
		msg = fmt.Sprintf("%d+ (use --users to get total count - may take awhile)", o.PublicMemberCount)
	} else {
		msg = fmt.Sprintf("%d", o.PublicMemberCount)
	}
	fmt.Println(aurora.Bold(aurora.Green("- Public Members:")), aurora.Bold(msg))
	fmt.Println(aurora.Bold(aurora.Green("- Followers:")), aurora.Bold(o.Followers))
	fmt.Println(aurora.Bold(aurora.Green("- Following:")), aurora.Bold(o.Following))
}

func getOrgByName(orgName string) *github.Organization {
	client := github.NewClient(nil)

	githubOrg, _, err := client.Organizations.Get(context.Background(), orgName)

	if err != nil {
		return nil
	}

	return githubOrg
}

func getOrgMembers(orgName string) []*github.User {
	client := github.NewClient(nil)

	opt := &github.ListMembersOptions{
		PublicOnly: true,
		ListOptions: github.ListOptions{
			Page:    1,
			PerPage: 100,
		}}

	// Get total org members
	orgMembers, resp, err := client.Organizations.ListMembers(context.Background(), orgName, opt)

	if err != nil {
		return nil
	}

	// If we have less than 100 users, just return what we have - otherwise move on to pagination requests..
	if len(orgMembers) < 100 || Users == false {
		return orgMembers
	}

	// Create a new temp variable to hold all the users in
	ret := []*github.User{}

	// add all the users we got to the return variable
	for _, member := range orgMembers {
		ret = append(ret, member)
	}

	// Run through every page, gathering all the members and adding them to the return value
	for i := 2; i < resp.LastPage+1; i++ {
		opt := &github.ListMembersOptions{
			PublicOnly: true,
			ListOptions: github.ListOptions{
				Page:    i,
				PerPage: 100,
			}}

		orgMembers, _, err := client.Organizations.ListMembers(context.Background(), orgName, opt)

		if err != nil {
			return nil
		}

		for _, member := range orgMembers {
			ret = append(ret, member)
		}
	}

	return ret
}
