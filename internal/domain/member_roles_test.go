package domain

import "testing"

func TestMemberRole_String(t *testing.T) {
	tests := []struct {
		name string
		role MemberRole
		want string
	}{
		{
			name: "default role",
			role: MemberRoleDefault,
			want: "default",
		},
		{
			name: "author role",
			role: MemberRolePrAuthor,
			want: "author",
		},
		{
			name: "reassigned role",
			role: MemberRoleHadReasigned,
			want: "reassigned",
		},
		{
			name: "unknown role",
			role: MemberRole(999),
			want: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.role.String(); got != tt.want {
				t.Errorf("MemberRole.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemberRoleFromString(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want MemberRole
	}{
		{
			name: "default string",
			s:    "default",
			want: MemberRoleDefault,
		},
		{
			name: "author string",
			s:    "author",
			want: MemberRolePrAuthor,
		},
		{
			name: "reassigned string",
			s:    "reassigned",
			want: MemberRoleHadReasigned,
		},
		{
			name: "unknown string",
			s:    "unknown",
			want: MemberStatusDefault,
		},
		{
			name: "empty string",
			s:    "",
			want: MemberStatusDefault,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MemberRoleFromString(tt.s); got != tt.want {
				t.Errorf("MemberRoleFromString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMembersRolesFromSliceOfStrings(t *testing.T) {
	tests := []struct {
		name string
		slS  []string
		want []MemberRole
	}{
		{
			name: "single role",
			slS:  []string{"default"},
			want: []MemberRole{MemberRoleDefault},
		},
		{
			name: "multiple roles",
			slS:  []string{"default", "author", "reassigned"},
			want: []MemberRole{MemberRoleDefault, MemberRolePrAuthor, MemberRoleHadReasigned},
		},
		{
			name: "empty slice",
			slS:  []string{},
			want: []MemberRole{},
		},
		{
			name: "unknown roles",
			slS:  []string{"unknown", "test"},
			want: []MemberRole{MemberStatusDefault, MemberStatusDefault},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MembersRolesFromSliceOfStrings(tt.slS)
			if len(got) != len(tt.want) {
				t.Errorf("MembersRolesFromSliceOfStrings() length = %v, want %v", len(got), len(tt.want))
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("MembersRolesFromSliceOfStrings()[%d] = %v, want %v", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestMemberWithRole(t *testing.T) {
	tests := []struct {
		name     string
		memberId MemberId
		role     MemberRole
	}{
		{
			name:     "member with default role",
			memberId: MemberId("member-1"),
			role:     MemberRoleDefault,
		},
		{
			name:     "member with author role",
			memberId: MemberId("member-2"),
			role:     MemberRolePrAuthor,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			memberWithRole := MemberWithRole{
				MemberID: tt.memberId,
				Role:     tt.role,
			}
			if memberWithRole.MemberID != tt.memberId {
				t.Errorf("MemberWithRole.MemberID = %v, want %v", memberWithRole.MemberID, tt.memberId)
			}
			if memberWithRole.Role != tt.role {
				t.Errorf("MemberWithRole.Role = %v, want %v", memberWithRole.Role, tt.role)
			}
		})
	}
}
