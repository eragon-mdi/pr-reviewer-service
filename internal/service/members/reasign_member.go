package servmembers

import (
	"context"
	"math/rand"

	"github.com/eragon-mdi/pr-reviewer-service/internal/domain"
)

func (ms *MembersService) ReasignMember(
	_ context.Context,
	memId domain.MemberId,
	mems domain.MembersHistories,
) (nilId domain.MemberId, _ error) {

	slMems := mems.Slice()
	l := len(slMems)

	start := rand.Intn(l)

	for i := range l {
		idx := (start + i) % l
		member := slMems[idx]

		if ms.allowedRoles.IsMemberAllowed(member) {
			return member.Id, nil
		}
	}

	return nilId, domain.ErrForbidden
}
