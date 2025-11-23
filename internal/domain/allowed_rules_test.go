package domain

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewAllowedRules(t *testing.T) {
	tests := []struct {
		name         string
		historyReuse bool
		statuses     []MemberStatus
		roles        []MemberRole
		wantStatuses int
		wantRoles    int
	}{
		{
			name:         "empty rules",
			historyReuse: false,
			statuses:     []MemberStatus{},
			roles:        []MemberRole{},
			wantStatuses: 0,
			wantRoles:    0,
		},
		{
			name:         "single status and role",
			historyReuse: true,
			statuses:     []MemberStatus{MemberStatusActive},
			roles:        []MemberRole{MemberRoleDefault},
			wantStatuses: 1,
			wantRoles:    1,
		},
		{
			name:         "multiple statuses and roles",
			historyReuse: false,
			statuses:     []MemberStatus{MemberStatusActive, MemberStatusInactive},
			roles:        []MemberRole{MemberRoleDefault, MemberRolePrAuthor},
			wantStatuses: 2,
			wantRoles:    2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ar := NewAllowedRules(tt.historyReuse, tt.statuses, tt.roles)
			if len(ar.allowedStatuses) != tt.wantStatuses {
				t.Errorf("AllowedRules.allowedStatuses length = %v, want %v", len(ar.allowedStatuses), tt.wantStatuses)
			}
			if len(ar.allowedRoles) != tt.wantRoles {
				t.Errorf("AllowedRules.allowedRoles length = %v, want %v", len(ar.allowedRoles), tt.wantRoles)
			}
		})
	}
}

func TestAllowedRules_IsMemberAllowed(t *testing.T) {
	tests := []struct {
		name            string
		historyReuse    bool
		allowedStatuses []MemberStatus
		allowedRoles    []MemberRole
		member          MemberHistory
		want            bool
	}{
		{
			name:            "allowed member - active, default role, history reuse true",
			historyReuse:    true,
			allowedStatuses: []MemberStatus{MemberStatusActive},
			allowedRoles:    []MemberRole{MemberRoleDefault},
			member: MemberHistory{
				Id:                MemberId(uuid.New().String()),
				Status:            MemberStatusActive,
				Role:              MemberRoleDefault,
				wasAssignedBefore: false,
			},
			want: true,
		},
		{
			name:            "not allowed - was assigned before, history reuse false",
			historyReuse:    false,
			allowedStatuses: []MemberStatus{MemberStatusActive},
			allowedRoles:    []MemberRole{MemberRoleDefault},
			member: MemberHistory{
				Id:                MemberId(uuid.New().String()),
				Status:            MemberStatusActive,
				Role:              MemberRoleDefault,
				wasAssignedBefore: true,
			},
			want: false,
		},
		{
			name:            "allowed - was assigned before, history reuse true",
			historyReuse:    true,
			allowedStatuses: []MemberStatus{MemberStatusActive},
			allowedRoles:    []MemberRole{MemberRoleDefault},
			member: MemberHistory{
				Id:                MemberId(uuid.New().String()),
				Status:            MemberStatusActive,
				Role:              MemberRoleDefault,
				wasAssignedBefore: true,
			},
			want: true,
		},
		{
			name:            "not allowed - wrong status",
			historyReuse:    true,
			allowedStatuses: []MemberStatus{MemberStatusActive},
			allowedRoles:    []MemberRole{MemberRoleDefault},
			member: MemberHistory{
				Id:                MemberId(uuid.New().String()),
				Status:            MemberStatusInactive,
				Role:              MemberRoleDefault,
				wasAssignedBefore: false,
			},
			want: false,
		},
		{
			name:            "not allowed - wrong role",
			historyReuse:    true,
			allowedStatuses: []MemberStatus{MemberStatusActive},
			allowedRoles:    []MemberRole{MemberRoleDefault},
			member: MemberHistory{
				Id:                MemberId(uuid.New().String()),
				Status:            MemberStatusActive,
				Role:              MemberRolePrAuthor,
				wasAssignedBefore: false,
			},
			want: false,
		},
		{
			name:            "not allowed - inactive status",
			historyReuse:    true,
			allowedStatuses: []MemberStatus{MemberStatusActive},
			allowedRoles:    []MemberRole{MemberRoleDefault},
			member: MemberHistory{
				Id:                MemberId(uuid.New().String()),
				Status:            MemberStatusInactive,
				Role:              MemberRoleDefault,
				wasAssignedBefore: false,
			},
			want: false,
		},
		{
			name:            "allowed - multiple statuses, one matches",
			historyReuse:    true,
			allowedStatuses: []MemberStatus{MemberStatusActive, MemberStatusInactive},
			allowedRoles:    []MemberRole{MemberRoleDefault},
			member: MemberHistory{
				Id:                MemberId(uuid.New().String()),
				Status:            MemberStatusInactive,
				Role:              MemberRoleDefault,
				wasAssignedBefore: false,
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ar := NewAllowedRules(tt.historyReuse, tt.allowedStatuses, tt.allowedRoles)
			if got := ar.IsMemberAllowed(tt.member); got != tt.want {
				t.Errorf("AllowedRules.IsMemberAllowed() = %v, want %v", got, tt.want)
			}
		})
	}
}
