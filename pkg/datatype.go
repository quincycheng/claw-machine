package pkg

type User struct {
	FirstName     string
	LastName      string
	MiddleName    string
	Email         string
	UserName      string
	ID            string
	RawUserDetail string
	RawUser       string
	Action        string
	ActionRemark  string
	NewName       string
	NewGroups     []string
}

type Group struct {
	GroupName      string
	GroupType      string
	Description    string
	ID             string
	RawGroup       string
	RawGroupDetail string
	Action         string
	ActionRemark   string
	NewRoleName    string
}
