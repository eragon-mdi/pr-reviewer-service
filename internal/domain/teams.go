package domain

type TeamName string

type Team struct {
	Name    TeamName
	Members Members
}

func NewTeam(name TeamName, members ...Member) Team {
	return Team{
		Name:    name,
		Members: members,
	}
}

func (tm TeamName) String() string {
	return string(tm)
}
