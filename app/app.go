package app

import (
	"fmt"
	"os"

	"github.com/erikgeiser/promptkit/selection"
	"github.com/quincycheng/claw-machine/pkg"
	"github.com/quincycheng/claw-machine/util"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/tidwall/gjson"
)

func Run(config util.Config) {

	util.PrintHeader("Claw Machine - migrating from Privilege Cloud to Identity Security Platform")

	s := util.StartSpinner("Loading configuration", "Configuration loaded")
	//DisplayConfig(config)
	s.Stop()
	fmt.Println()

	config.From.Password = util.InputPassword("Please input password of " + config.From.User + " from " + config.From.Url)
	config.To.Password = util.InputPassword("Please input password of " + config.To.User + " from " + config.To.Url)

	s = util.StartSpinner("Connecting to Privilege Cloud Standalone", "Connected to Privilege Cloud Standalone ")
	fromToken, err := pkg.CyberArkAuthnStandalone(config)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	s.Stop()
	fmt.Println()

	s = util.StartSpinner("Getting Groups from Privilege Cloud Standalone", "Groups from Privilege Cloud Standalone")
	allExistingGroups, err := GetAllGroupsStandalone(config, fromToken)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	s.Stop()
	fmt.Println()

	s = util.StartSpinner("Getting Users from Privilege Cloud Standalone", "Users from Privilege Cloud Standalone")
	allExistingUsers, err := GetAllUsersStandalone(config, fromToken)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	s.Stop()
	fmt.Println()

	s = util.StartSpinner("Connecting to Identity Security Platform", "Connected to Identity Security Platform")
	toToken, err := pkg.AuthnSharedService(config)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	s.Stop()
	fmt.Println()

	s = util.StartSpinner("Getting Groups from Identity Security Platform", "Groups from Identity Security Platform")
	allISPRole, err := pkg.GetAllRolesISP(config, toToken)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	s.Stop()
	fmt.Println()

	s = util.StartSpinner("Getting Users from Identity Security Platform", "Users from Identity Security Platform")
	allISPUser, err := pkg.GetAllUsersISP(config, toToken)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	s.Stop()
	fmt.Println()

	// TODO: Processing Groups
	allExistingGroups = ProcessGroups(config, allExistingGroups, allISPRole)

	// TODO: Processing Users
	// TODO: Show Options - Review & Migrate

	const ReviewConfig = "Review Configuration"

	const ReviewGroupList = "Review Group List from Privilege Cloud"
	const ReviewUserList = "Review User List from Privilege Cloud"
	const ReviewUserDetail = "Review Specific User Details from Privilege Cloud"

	const ReviewISPUserList = "Review User List from Identity Security Platform (json)"
	const ReviewISPRoleList = "Review Role List from Identity Security Platform (json)"

	const ShowDebug = "Show Debug details for troubleshooting"
	const Migrate = "Migrate Users & Groups"
	const Quit = "Quit"

	for {
		fmt.Println()
		sp := selection.New("Please select your action: ", selection.Choices([]string{ReviewConfig, ReviewUserList, ReviewUserDetail,
			ReviewGroupList, ReviewISPUserList, ReviewISPRoleList, ShowDebug, Migrate, Quit}))

		choice, err := sp.RunPrompt()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		switch choice.Value {
		case Quit:
			fmt.Println("Bye!")
			os.Exit(0)
		case ReviewConfig:
			util.PrintHeader(ReviewConfig)
			DisplayConfig(config)
		case ReviewUserList:
			DisplayAllUsersStandalone(allExistingUsers)
		case ReviewGroupList:
			DisplayAllGroupsStandalone(allExistingGroups)
		case ReviewUserDetail:
			DisplaySpecificUserDetails(allExistingUsers)
		case ReviewISPUserList:
			util.PrintHeader(ReviewConfig)
			util.PrintPrettyJson(allISPUser)
		case ReviewISPRoleList:
			util.PrintHeader(ReviewConfig)
			util.PrintPrettyJson(allISPRole)
		case ShowDebug:
			util.PrintHeader(ReviewConfig)
			DisplayConfig(config)
			util.PrintHeader(ReviewUserList)
			DisplayAllUsersStandalone(allExistingUsers)
			util.PrintHeader(ReviewGroupList)
			DisplayAllGroupsStandalone(allExistingGroups)
			util.PrintHeader(ReviewISPUserList)
			util.PrintPrettyJson(allISPUser)
			util.PrintHeader(ReviewISPUserList)
			util.PrintPrettyJson(allISPRole)
		}
		_ = choice
	}

	_ = allExistingGroups
	_ = allExistingUsers
}

func DisplayConfig(config util.Config) {
	fmt.Println("From: ", config.From.Url)
	fmt.Println("To: ", config.To.Url)

	fmt.Println("Group Mapping:")
	for _, s := range config.Mapping.Groups {
		fmt.Println(" - ", s.From, ">>", s.To)
	}

	fmt.Println("Role Mappings:")
	for _, s := range config.Mapping.Roles {
		fmt.Println(" - ", s.From, ">>", s.To)
	}

	fmt.Println("Skip Following Users:")
	for _, s := range config.Skip.Users {
		fmt.Println(" - ", s)
	}

	fmt.Println("Skip Following Groups:")
	for _, s := range config.Skip.Groups {
		fmt.Println(" - ", s)

	}
}

func DisplaySpecificUserDetails(users []pkg.User) {

	var choiceList []selection.Choice
	for _, user := range users {
		displayString := user.ID + ". " + user.UserName
		if user.Email != "" {
			displayString += " <" + user.Email + ">"
		}
		choiceList = append(choiceList, selection.Choice{String: displayString, Value: user.ID})
	}

	sp := selection.New("Please select user for review: ", selection.Choices(choiceList))

	choice, err := sp.RunPrompt()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	//	switch choice.Value {
	//	case Quit:
	//		fmt.Println("Bye!")
	//		os.Exit(0)
	//	}
	fmt.Println(choice.Value)
}

func DisplayAllUsersStandalone(users []pkg.User) {
	fmt.Print("Total no of Users: ")
	fmt.Println(len(users))

	t := table.NewWriter()
	t.SetStyle(table.StyleLight)
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"ID", "User Name", "First Name", "Last Name", "Email", "Action", "Remarks"})

	for _, user := range users {
		t.AppendRow([]interface{}{user.ID, user.UserName, user.FirstName, user.LastName, user.Email, user.Action, user.ActionRemark})
	}
	t.Render()

}

func DisplayAllGroupsStandalone(groups []pkg.Group) {
	fmt.Print("Total no of groups: ")
	fmt.Println(len(groups))

	t := table.NewWriter()
	t.SetStyle(table.StyleLight)
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"ID", "Name", "Type", "Description", "Action", "Remarks", "New Role Name"})

	for _, group := range groups {
		t.AppendRow([]interface{}{group.ID, group.GroupName, group.GroupType, group.Description, group.Action, group.ActionRemark, group.NewRoleName})
	}
	t.Render()
}

func GetAllGroupsStandalone(config util.Config, token string) ([]pkg.Group, error) {
	var result []pkg.Group
	theUrl := "INIT"
	for ok := true; ok; ok = (theUrl != "") {

		jsonString, _ := pkg.GetAllGroupsStandalone(config, token, theUrl)

		for _, val := range gjson.Get(jsonString, "value").Array() {
			tmpGroup := pkg.Group{
				GroupName:   val.Get("groupName").String(),
				GroupType:   val.Get("groupType").String(),
				Description: val.Get("description").String(),
				ID:          val.Get("id").String(),
				RawGroup:    val.String(),
			}

			jsonGroupDetail, _ := pkg.GetGroupDetailStandalone(config, token, tmpGroup.ID)
			tmpGroup.RawGroupDetail = jsonGroupDetail
			result = append(result, tmpGroup)
		}
		theUrl = gjson.Get(jsonString, "nextLink").String()
	}
	return result, nil
}

func GetAllUsersStandalone(config util.Config, token string) ([]pkg.User, error) {
	var result []pkg.User
	theUrl := "INIT"
	for ok := true; ok; ok = (theUrl != "") {
		jsonString, _ := pkg.GetAllUsersStandalone(config, token, theUrl)

		for _, val := range gjson.Get(jsonString, "Users").Array() {
			tmpUser := pkg.User{
				ID:         val.Get("id").String(),
				FirstName:  val.Get("firstname").String(),
				MiddleName: val.Get("middlename").String(),
				LastName:   val.Get("lastname").String(),
				UserName:   val.Get("username").String(),
				RawUser:    val.Str,
			}

			jsonUserDetail, _ := pkg.GetUserDetailStandalone(config, token, tmpUser.ID)
			tmpUser.RawUserDetail = jsonUserDetail
			tmpUser.Email = gjson.Get(jsonUserDetail, "internet.businessEmail").String()

			result = append(result, tmpUser)
			//fmt.Println()
			//fmt.Println(tmpUser.ID + ") " + tmpUser.UserName + ": " + tmpUser.Email)
		}
		theUrl = gjson.Get(jsonString, "nextLink").String()
	}
	return result, nil
}

func ProcessGroups(config util.Config, groups []pkg.Group, ispRolesJson string) []pkg.Group {

	for i, group := range groups {

		// Convert groups to roles based on config
		for _, s := range config.Mapping.Groups {
			if group.GroupName == s.From {
				group.Action = "Migrate"
				group.ActionRemark = "Convert based on config"
				group.NewRoleName = s.To
			}
		}
		// Convert groups to roles based on config
		if group.Action == "" {
			for _, s := range config.Skip.Groups {
				if group.GroupName == s {
					group.Action = "Skip"
					group.ActionRemark = "Skipped based on config"
				}
			}
		}

		// check if role exist based on original group name
		if group.Action == "" {
			result := gjson.Get(ispRolesJson, "Result.Results.#.Row.Name")
			for _, name := range result.Array() {
				//fmt.Println(name.String())
				if group.GroupName == name.String() {
					group.Action = "Skip"
					group.ActionRemark = "Skipped as Role already exist"
				}
			}
		}

		// check if role exist based on new converted group name
		if group.ActionRemark == "Convert based on config" {
			result := gjson.Get(ispRolesJson, "Result.Results.#.Row.Name")
			for _, name := range result.Array() {
				//fmt.Println(name.String())
				if group.NewRoleName == name.String() {
					group.Action = "Skip"
					group.ActionRemark = "Skipped as Converted Role already exist"
				}
			}
		}

		if group.Action == "" {
			group.Action = "Migrate"
			group.NewRoleName = group.GroupName
		}

		groups[i] = group
	}
	return groups
}
