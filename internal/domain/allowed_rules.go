package domain

type AllowedRules struct {
	allowedStatuses map[MemberStatus]struct{}
	allowedRoles    map[MemberRole]struct{}
	historyReuse    bool
}

func NewAllowedRules(historyReuse bool, statuses []MemberStatus, roles []MemberRole) AllowedRules {
	mS := make(map[MemberStatus]struct{}, len(statuses))
	for _, s := range statuses {
		mS[s] = struct{}{}
	}

	mR := make(map[MemberRole]struct{}, len(roles))
	for _, r := range roles {
		mR[r] = struct{}{}
	}

	return AllowedRules{
		allowedStatuses: mS,
		allowedRoles:    mR,
		historyReuse:    historyReuse,
	}
}

func (ar AllowedRules) IsMemberAllowed(member MemberHistory) bool {
	if !ar.historyReuse && member.wasAssignedBefore {
		return false
	}

	_, ok1 := ar.allowedStatuses[member.Status]

	_, ok2 := ar.allowedRoles[member.Role]

	return ok1 && ok2
}
