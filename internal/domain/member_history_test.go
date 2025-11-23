package domain

import (
	"testing"

	"github.com/google/uuid"
)

func TestMembersHistories_Slice(t *testing.T) {
	tests := []struct {
		name string
		mh   MembersHistories
		want int
	}{
		{
			name: "empty histories",
			mh:   MembersHistories{},
			want: 0,
		},
		{
			name: "single history",
			mh: MembersHistories{
				{
					Id:                MemberId(uuid.New().String()),
					Status:            MemberStatusActive,
					Role:              MemberRoleDefault,
					wasAssignedBefore: false,
				},
			},
			want: 1,
		},
		{
			name: "multiple histories",
			mh: MembersHistories{
				{
					Id:                MemberId(uuid.New().String()),
					Status:            MemberStatusActive,
					Role:              MemberRoleDefault,
					wasAssignedBefore: false,
				},
				{
					Id:                MemberId(uuid.New().String()),
					Status:            MemberStatusInactive,
					Role:              MemberRolePrAuthor,
					wasAssignedBefore: true,
				},
			},
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.mh.Slice()
			if len(got) != tt.want {
				t.Errorf("MembersHistories.Slice() length = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func TestMembersHistories_Empty(t *testing.T) {
	tests := []struct {
		name string
		mh   MembersHistories
		want bool
	}{
		{
			name: "empty histories",
			mh:   MembersHistories{},
			want: true,
		},
		{
			name: "nil histories",
			mh:   nil,
			want: true,
		},
		{
			name: "non-empty histories",
			mh: MembersHistories{
				{
					Id:     MemberId(uuid.New().String()),
					Status: MemberStatusActive,
					Role:   MemberRoleDefault,
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.mh.Empty(); got != tt.want {
				t.Errorf("MembersHistories.Empty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemberHistory(t *testing.T) {
	tests := []struct {
		name              string
		id                MemberId
		status            MemberStatus
		role              MemberRole
		wasAssignedBefore bool
	}{
		{
			name:              "active member not assigned before",
			id:                MemberId(uuid.New().String()),
			status:            MemberStatusActive,
			role:              MemberRoleDefault,
			wasAssignedBefore: false,
		},
		{
			name:              "inactive member assigned before",
			id:                MemberId(uuid.New().String()),
			status:            MemberStatusInactive,
			role:              MemberRoleHadReasigned,
			wasAssignedBefore: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mh := MemberHistory{
				Id:                tt.id,
				Status:            tt.status,
				Role:              tt.role,
				wasAssignedBefore: tt.wasAssignedBefore,
			}
			if mh.Id != tt.id {
				t.Errorf("MemberHistory.Id = %v, want %v", mh.Id, tt.id)
			}
			if mh.Status != tt.status {
				t.Errorf("MemberHistory.Status = %v, want %v", mh.Status, tt.status)
			}
			if mh.Role != tt.role {
				t.Errorf("MemberHistory.Role = %v, want %v", mh.Role, tt.role)
			}
			if mh.wasAssignedBefore != tt.wasAssignedBefore {
				t.Errorf("MemberHistory.wasAssignedBefore = %v, want %v", mh.wasAssignedBefore, tt.wasAssignedBefore)
			}
		})
	}
}
