package syncer

type Syncer interface {
	SyncUserRanks([]Member) []Member
}

type RankMapping struct {
	RoleId   string
	RoleName string
}

type Member struct {
	DiscId        string
	UserName      string
	DiscRankId    string
	DiscRankName  string
	Roles         []string
	RolesToAdd    []string
	RolesToRemove []string
}
