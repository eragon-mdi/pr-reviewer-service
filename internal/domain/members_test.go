package domain

import (
	"testing"

	"github.com/google/uuid"
)

func TestMemberId_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		memberId MemberId
		want     bool
	}{
		{
			name:     "valid uuid",
			memberId: MemberId(uuid.New().String()),
			want:     true,
		},
		{
			name:     "invalid uuid - empty string",
			memberId: MemberId(""),
			want:     false,
		},
		{
			name:     "invalid uuid - random string",
			memberId: MemberId("not-a-uuid"),
			want:     false,
		},
		{
			name:     "invalid uuid - number",
			memberId: MemberId("12345"),
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.memberId.IsValid(); got != tt.want {
				t.Errorf("MemberId.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemberId_String(t *testing.T) {
	testUUID := uuid.New().String()

	tests := []struct {
		name     string
		memberId MemberId
		want     string
	}{
		{
			name:     "simple string",
			memberId: MemberId("test-id"),
			want:     "test-id",
		},
		{
			name:     "uuid string",
			memberId: MemberId(testUUID),
			want:     testUUID,
		},
		{
			name:     "empty string",
			memberId: MemberId(""),
			want:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.memberId.String(); got != tt.want {
				t.Errorf("MemberId.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemberStatus_String(t *testing.T) {
	tests := []struct {
		name   string
		status MemberStatus
		want   string
	}{
		{
			name:   "active status",
			status: MemberStatusActive,
			want:   "active",
		},
		{
			name:   "inactive status",
			status: MemberStatusInactive,
			want:   "inactive",
		},
		{
			name:   "unknown status",
			status: MemberStatus(999),
			want:   "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.status.String(); got != tt.want {
				t.Errorf("MemberStatus.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemberStatus_IsActive(t *testing.T) {
	tests := []struct {
		name   string
		status MemberStatus
		want   bool
	}{
		{
			name:   "active status",
			status: MemberStatusActive,
			want:   true,
		},
		{
			name:   "inactive status",
			status: MemberStatusInactive,
			want:   false,
		},
		{
			name:   "unknown status",
			status: MemberStatus(999),
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.status.IsActive(); got != tt.want {
				t.Errorf("MemberStatus.IsActive() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemberStatusFromString(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want MemberStatus
	}{
		{
			name: "active string",
			s:    "active",
			want: MemberStatusActive,
		},
		{
			name: "inactive string",
			s:    "inactive",
			want: MemberStatusInactive,
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
			if got := MemberStatusFromString(tt.s); got != tt.want {
				t.Errorf("MemberStatusFromString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemberStatusIsActiveByBool(t *testing.T) {
	tests := []struct {
		name     string
		isActive bool
		want     MemberStatus
	}{
		{
			name:     "true to active",
			isActive: true,
			want:     MemberStatusActive,
		},
		{
			name:     "false to inactive",
			isActive: false,
			want:     MemberStatusInactive,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MemberStatusIsActiveByBool(tt.isActive); got != tt.want {
				t.Errorf("MemberStatusIsActiveByBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMembersStatusesFromSliceOfStrings(t *testing.T) {
	tests := []struct {
		name string
		slS  []string
		want []MemberStatus
	}{
		{
			name: "single active",
			slS:  []string{"active"},
			want: []MemberStatus{MemberStatusActive},
		},
		{
			name: "mixed statuses",
			slS:  []string{"active", "inactive"},
			want: []MemberStatus{MemberStatusActive, MemberStatusInactive},
		},
		{
			name: "empty slice",
			slS:  []string{},
			want: []MemberStatus{},
		},
		{
			name: "unknown statuses",
			slS:  []string{"unknown", "test"},
			want: []MemberStatus{MemberStatusDefault, MemberStatusDefault},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MembersStatusesFromSliceOfStrings(tt.slS)
			if len(got) != len(tt.want) {
				t.Errorf("MembersStatusesFromSliceOfStrings() length = %v, want %v", len(got), len(tt.want))
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("MembersStatusesFromSliceOfStrings()[%d] = %v, want %v", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestMemberBuilder(t *testing.T) {
	tests := []struct {
		name        string
		id          MemberId
		setup       func(builder memberBuilder) memberBuilder
		wantName    string
		wantStatus  MemberStatus
		wantReviews int
	}{
		{
			name: "minimal member",
			id:   MemberId(uuid.New().String()),
			setup: func(builder memberBuilder) memberBuilder {
				return builder
			},
			wantName:    "",
			wantStatus:  MemberStatusDefault,
			wantReviews: 0,
		},
		{
			name: "member with name",
			id:   MemberId(uuid.New().String()),
			setup: func(builder memberBuilder) memberBuilder {
				return builder.Name("Test User")
			},
			wantName:    "Test User",
			wantStatus:  MemberStatusDefault,
			wantReviews: 0,
		},
		{
			name: "member with status",
			id:   MemberId(uuid.New().String()),
			setup: func(builder memberBuilder) memberBuilder {
				return builder.Status(MemberStatusActive)
			},
			wantName:    "",
			wantStatus:  MemberStatusActive,
			wantReviews: 0,
		},
		{
			name: "full member",
			id:   MemberId(uuid.New().String()),
			setup: func(builder memberBuilder) memberBuilder {
				return builder.
					Name("Test User").
					Status(MemberStatusActive).
					Reviews([]PullRequestShort{
						{Id: "pr1", Name: "PR1", AuthorId: MemberId(uuid.New().String()), Status: PrStatusOpen},
					})
			},
			wantName:    "Test User",
			wantStatus:  MemberStatusActive,
			wantReviews: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := MemberBuilder(tt.id)
			builder = tt.setup(builder)
			member := builder.Build()

			if member.Id != tt.id {
				t.Errorf("Member.Id = %v, want %v", member.Id, tt.id)
			}
			if member.Name != tt.wantName {
				t.Errorf("Member.Name = %v, want %v", member.Name, tt.wantName)
			}
			if member.Status != tt.wantStatus {
				t.Errorf("Member.Status = %v, want %v", member.Status, tt.wantStatus)
			}
			if len(member.Reviews) != tt.wantReviews {
				t.Errorf("Member.Reviews length = %v, want %v", len(member.Reviews), tt.wantReviews)
			}
		})
	}
}

func TestMembers_Empty(t *testing.T) {
	tests := []struct {
		name string
		m    Members
		want bool
	}{
		{
			name: "empty members",
			m:    Members{},
			want: true,
		},
		{
			name: "nil members",
			m:    nil,
			want: true,
		},
		{
			name: "non-empty members",
			m: Members{
				{Id: MemberId(uuid.New().String()), Name: "User1"},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Empty(); got != tt.want {
				t.Errorf("Members.Empty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMembers_Slice(t *testing.T) {
	tests := []struct {
		name string
		m    Members
		want int
	}{
		{
			name: "empty slice",
			m:    Members{},
			want: 0,
		},
		{
			name: "single member",
			m: Members{
				{Id: MemberId(uuid.New().String()), Name: "User1"},
			},
			want: 1,
		},
		{
			name: "multiple members",
			m: Members{
				{Id: MemberId(uuid.New().String()), Name: "User1"},
				{Id: MemberId(uuid.New().String()), Name: "User2"},
			},
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.m.Slice()
			if len(got) != tt.want {
				t.Errorf("Members.Slice() length = %v, want %v", len(got), tt.want)
			}
		})
	}
}
