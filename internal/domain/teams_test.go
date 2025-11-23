package domain

import (
	"testing"

	"github.com/google/uuid"
)

func TestTeamName_String(t *testing.T) {
	tests := []struct {
		name string
		tn   TeamName
		want string
	}{
		{
			name: "simple team name",
			tn:   TeamName("backend"),
			want: "backend",
		},
		{
			name: "empty team name",
			tn:   TeamName(""),
			want: "",
		},
		{
			name: "team name with spaces",
			tn:   TeamName("frontend team"),
			want: "frontend team",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tn.String(); got != tt.want {
				t.Errorf("TeamName.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewTeam(t *testing.T) {
	tests := []struct {
		name     string
		teamName TeamName
		members  []Member
		wantName TeamName
		wantLen  int
	}{
		{
			name:     "team without members",
			teamName: TeamName("backend"),
			members:  []Member{},
			wantName: TeamName("backend"),
			wantLen:  0,
		},
		{
			name:     "team with single member",
			teamName: TeamName("backend"),
			members: []Member{
				{Id: MemberId(uuid.New().String()), Name: "User1", Status: MemberStatusActive},
			},
			wantName: TeamName("backend"),
			wantLen:  1,
		},
		{
			name:     "team with multiple members",
			teamName: TeamName("frontend"),
			members: []Member{
				{Id: MemberId(uuid.New().String()), Name: "User1", Status: MemberStatusActive},
				{Id: MemberId(uuid.New().String()), Name: "User2", Status: MemberStatusActive},
				{Id: MemberId(uuid.New().String()), Name: "User3", Status: MemberStatusInactive},
			},
			wantName: TeamName("frontend"),
			wantLen:  3,
		},
		{
			name:     "team with variadic members",
			teamName: TeamName("devops"),
			members: []Member{
				{Id: MemberId(uuid.New().String()), Name: "User1"},
			},
			wantName: TeamName("devops"),
			wantLen:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			team := NewTeam(tt.teamName, tt.members...)
			if team.Name != tt.wantName {
				t.Errorf("Team.Name = %v, want %v", team.Name, tt.wantName)
			}
			if len(team.Members) != tt.wantLen {
				t.Errorf("Team.Members length = %v, want %v", len(team.Members), tt.wantLen)
			}
		})
	}
}
