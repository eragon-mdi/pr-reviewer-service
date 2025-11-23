package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestPrId_String(t *testing.T) {
	testUUID := uuid.New().String()

	tests := []struct {
		name string
		pId  PrId
		want string
	}{
		{
			name: "simple pr id",
			pId:  PrId("pr-123"),
			want: "pr-123",
		},
		{
			name: "uuid pr id",
			pId:  PrId(testUUID),
			want: testUUID,
		},
		{
			name: "empty pr id",
			pId:  PrId(""),
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.pId.String(); got != tt.want {
				t.Errorf("PrId.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrName_String(t *testing.T) {
	tests := []struct {
		name  string
		pName PrName
		want  string
	}{
		{
			name:  "simple pr name",
			pName: PrName("Add feature"),
			want:  "Add feature",
		},
		{
			name:  "empty pr name",
			pName: PrName(""),
			want:  "",
		},
		{
			name:  "long pr name",
			pName: PrName("Implement complex feature with many changes"),
			want:  "Implement complex feature with many changes",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.pName.String(); got != tt.want {
				t.Errorf("PrName.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPullRequestShort_Create(t *testing.T) {
	tests := []struct {
		name  string
		prs   PullRequestShort
		check func(t *testing.T, pr PullRequest)
	}{
		{
			name: "create from short",
			prs: PullRequestShort{
				Id:       PrId("pr-123"),
				Name:     PrName("Test PR"),
				AuthorId: MemberId(uuid.New().String()),
				Status:   PrStatusOpen,
			},
			check: func(t *testing.T, pr PullRequest) {
				if pr.Id != PrId("pr-123") {
					t.Errorf("PullRequest.Id = %v, want pr-123", pr.Id)
				}
				if pr.Name != PrName("Test PR") {
					t.Errorf("PullRequest.Name = %v, want Test PR", pr.Name)
				}
				if pr.Status != PrStatusDefault {
					t.Errorf("PullRequest.Status = %v, want %v", pr.Status, PrStatusDefault)
				}
				if pr.CreatedAt.IsZero() {
					t.Error("PullRequest.CreatedAt should not be zero")
				}
				if !pr.MergedAt.IsZero() {
					t.Error("PullRequest.MergedAt should be zero")
				}
				if pr.version != 0 {
					t.Errorf("PullRequest.version = %v, want 0", pr.version)
				}
				if pr.AssignedReviews != nil {
					t.Error("PullRequest.AssignedReviews should be nil")
				}
			},
		},
		{
			name: "create with merged status",
			prs: PullRequestShort{
				Id:       PrId("pr-456"),
				Name:     PrName("Merged PR"),
				AuthorId: MemberId(uuid.New().String()),
				Status:   PrStatusMerged,
			},
			check: func(t *testing.T, pr PullRequest) {
				if pr.Status != PrStatusDefault {
					t.Errorf("PullRequest.Status should be reset to default, got %v", pr.Status)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pr := tt.prs.Create()
			tt.check(t, pr)
		})
	}
}

func TestPullRequests_Empty(t *testing.T) {
	tests := []struct {
		name string
		prs  PullRequests
		want bool
	}{
		{
			name: "empty pull requests",
			prs:  PullRequests{},
			want: true,
		},
		{
			name: "nil pull requests",
			prs:  nil,
			want: true,
		},
		{
			name: "non-empty pull requests",
			prs: PullRequests{
				{Id: PrId("pr-1"), Name: PrName("PR1"), AuthorId: MemberId(uuid.New().String()), Status: PrStatusOpen},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.prs.Empty(); got != tt.want {
				t.Errorf("PullRequests.Empty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrReasignMember(t *testing.T) {
	tests := []struct {
		name     string
		prId     PrId
		memberId MemberId
	}{
		{
			name:     "valid reassign",
			prId:     PrId("pr-123"),
			memberId: MemberId(uuid.New().String()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prReasMem := PrReasignMember{
				PrId:     tt.prId,
				MemberId: tt.memberId,
			}
			if prReasMem.PrId != tt.prId {
				t.Errorf("PrReasignMember.PrId = %v, want %v", prReasMem.PrId, tt.prId)
			}
			if prReasMem.MemberId != tt.memberId {
				t.Errorf("PrReasignMember.MemberId = %v, want %v", prReasMem.MemberId, tt.memberId)
			}
		})
	}
}

func TestPrWithReasignMember(t *testing.T) {
	tests := []struct {
		name string
		pr   PullRequest
		mId  MemberId
	}{
		{
			name: "pr with reassign member",
			pr: PullRequest{
				Id:        PrId("pr-123"),
				Name:      PrName("Test PR"),
				AuthorId:  MemberId(uuid.New().String()),
				Status:    PrStatusOpen,
				CreatedAt: time.Now(),
			},
			mId: MemberId(uuid.New().String()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prWithReasign := PrWithReasignMember{
				PullRequest: tt.pr,
				MemberId:    tt.mId,
			}
			if prWithReasign.PullRequest.Id != tt.pr.Id {
				t.Errorf("PrWithReasignMember.PullRequest.Id = %v, want %v", prWithReasign.PullRequest.Id, tt.pr.Id)
			}
			if prWithReasign.MemberId != tt.mId {
				t.Errorf("PrWithReasignMember.MemberId = %v, want %v", prWithReasign.MemberId, tt.mId)
			}
		})
	}
}
