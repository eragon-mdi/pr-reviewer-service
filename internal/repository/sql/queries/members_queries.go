package queries

const (
	UpdateMemberStatus = `
		UPDATE members
		SET is_active = $2
		WHERE uuid = $1
		RETURNING id, uuid, name, is_active;
	`

	GetPrReviewsByMember = `
		SELECT 
			pr.uuid AS pull_request_id,
			pr.title AS pull_request_name,
			author.uuid AS author_id,
			s.status AS status
		FROM pr_members pm
		INNER JOIN pull_requests pr ON pm.pr_id = pr.id
		INNER JOIN members reviewer ON pm.member_id = reviewer.id
		INNER JOIN members author ON pr.author_id = author.id
		INNER JOIN statuses s ON pr.status_id = s.id
		INNER JOIN roles r ON pm.role_id = r.id
		WHERE reviewer.uuid = $1
		  AND r.role = 'reviewer'
		ORDER BY pr.created_at DESC;
	`

	GetMemberByUUID = `
		SELECT id, uuid, name, is_active
		FROM members
		WHERE uuid = $1;
	`

	GetActiveMembersByTeamId = `
		SELECT m.id, m.uuid, m.name, m.is_active
		FROM members m
		INNER JOIN members_teams mt ON m.id = mt.member_id
		WHERE mt.team_id = $1
		  AND m.is_active = true
		  AND m.uuid != $2
		ORDER BY m.name;
	`
)
