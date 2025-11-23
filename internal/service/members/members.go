package servmembers

import (
	"fmt"

	"github.com/eragon-mdi/pr-reviewer-service/internal/domain"
	"github.com/go-faster/errors"
)

type MembersRepository interface {
	UpdateMemberStatus(domain.MemberId, domain.MemberStatus) (domain.Member, error)
	GetPrReviewsByMember(domain.MemberId) (domain.PullRequests, error)
}

func (ms *MembersService) SetMemberIsActive(member domain.Member) (domain.Member, error) {

	updMember, err := ms.repo.UpdateMemberStatus(member.Id, member.Status)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.Member{}, domain.ErrNotFound
		}
		return domain.Member{}, fmt.Errorf("%w: %w", domain.ErrInternal, err)
	}

	return updMember, nil
}

func (ms *MembersService) MemberReviews(id domain.MemberId) (domain.Member, error) {

	revs, err := ms.repo.GetPrReviewsByMember(id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.Member{}, domain.ErrNotFound
		}
		return domain.Member{}, fmt.Errorf("%w: %w", domain.ErrInternal, err)
	}

	if revs.Empty() {
		return domain.MemberBuilder(id).Build(), domain.ErrNoContent
	}

	return domain.MemberBuilder(id).Reviews(revs).Build(), nil
}
