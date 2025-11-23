package domain

type MemberRole int

const (
	MemberRoleDefault = iota
	MemberRolePrAuthor
	MemberRoleHadReasigned
)

type MemberWithRole struct {
	MemberID MemberId
	Role     MemberRole
}

var MemberRoleNames = map[MemberRole]string{
	MemberRoleDefault:      "default",
	MemberRolePrAuthor:     "author",
	MemberRoleHadReasigned: "reassigned",
}

func (r MemberRole) String() string {
	if name, ok := MemberRoleNames[r]; ok {
		return name
	}
	return "unknown"
}

func MemberRoleFromString(s string) MemberRole {
	for k, v := range MemberRoleNames {
		if v == s {
			return k
		}
	}
	return MemberStatusDefault
}

func MembersRolesFromSliceOfStrings(slS []string) []MemberRole {
	res := make([]MemberRole, 0, len(slS))

	for _, s := range slS {
		res = append(res, MemberRoleFromString(s))
	}

	return res
}
