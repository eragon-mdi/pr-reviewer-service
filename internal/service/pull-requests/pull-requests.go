package servpullrequests

import (
	"context"
	"fmt"

	"github.com/eragon-mdi/pr-reviewer-service/internal/domain"
	"github.com/go-faster/errors"
)

type PullRequestsRepository interface {
	CreatePullRequest(domain.PullRequest) (domain.PullRequest, error)
	MergePullRequest(domain.PrId) (domain.PullRequest, error)
	BeginReasignTx(context.Context) (ReassignTx, error)
}

type ReassignTx interface {
	GetPullRequestMembersHistories(context.Context, domain.PrId) (domain.MembersHistories, error)
	AssignMember(context.Context, domain.PrId, domain.MemberId) (domain.PullRequest, error)
	Commit() error
	Rollback() error
}

type MemberService interface {
	ReasignMember(context.Context, domain.MemberId, domain.MembersHistories) (domain.MemberId, error)
}

func (ps *PrService) Reasign(ctx context.Context, prReasMem domain.PrReasignMember) (_ domain.PrWithReasignMember, err error) {
	tx, err := ps.repo.BeginReasignTx(ctx)
	if err != nil {
		return domain.PrWithReasignMember{}, fmt.Errorf("%w: %w", domain.ErrInternal, err)
	}
	defer func() {
		if err == nil {
			if errCommit := tx.Commit(); errCommit != nil {
				err = fmt.Errorf("%w: %w", domain.ErrInternal, errCommit)
			}
			return
		}
		if errRollback := tx.Rollback(); errRollback != nil {
			err = fmt.Errorf("%w: %w", err, errRollback)
		}
	}()

	candidatesHistories, err := tx.GetPullRequestMembersHistories(ctx, prReasMem.PrId)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.PrWithReasignMember{}, domain.ErrNotFound
		}
		if errors.Is(err, domain.ErrNoContent) {
			return domain.PrWithReasignMember{}, domain.ErrNoContent
		}
		return domain.PrWithReasignMember{}, fmt.Errorf("%w: %w", domain.ErrInternal, err)
	}

	memberIdToAssign, err := ps.memServ.ReasignMember(ctx, prReasMem.MemberId, candidatesHistories)
	if err != nil {
		return domain.PrWithReasignMember{}, err
	}

	pr, err := tx.AssignMember(ctx, prReasMem.PrId, memberIdToAssign)
	if err != nil {
		return domain.PrWithReasignMember{}, fmt.Errorf("%w: %w", domain.ErrInternal, err)
	}

	return domain.PrWithReasignMember{
		PullRequest: pr,
		MemberId:    memberIdToAssign,
	}, nil
}

func (ps *PrService) NewPullRequest(basePR domain.PullRequestShort) (domain.PullRequest, error) {

	pr := basePR.Create()

	createdPr, err := ps.repo.CreatePullRequest(pr)
	if err != nil {
		if errors.Is(err, domain.ErrDuplicate) {
			return domain.PullRequest{}, domain.ErrDuplicate
		}
		if errors.Is(err, domain.ErrNotFound) {
			return domain.PullRequest{}, domain.ErrNotFound
		}
		return domain.PullRequest{}, fmt.Errorf("%w: %w", domain.ErrInternal, err)
	}

	return createdPr, nil
}

func (ps *PrService) Merge(id domain.PrId) (domain.PullRequest, error) {
	merged, err := ps.repo.MergePullRequest(id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.PullRequest{}, domain.ErrNotFound
		}
		if errors.Is(err, domain.ErrConflict) {
			return domain.PullRequest{}, domain.ErrConflict
		}
		return domain.PullRequest{}, fmt.Errorf("%w: %w", domain.ErrInternal, err)
	}

	return merged, nil
}
