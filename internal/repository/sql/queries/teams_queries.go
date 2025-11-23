package queries

const (
	CreateTeamWithMembers = `
		WITH team_ins AS (
			INSERT INTO teams (name)
			VALUES ($1)
			ON CONFLICT (name) DO NOTHING
			RETURNING id
		),
		team_sel AS (
			SELECT id FROM teams WHERE name = $1
		),
		team_id AS (
			SELECT id FROM team_ins
			UNION ALL
			SELECT id FROM team_sel
		)
		INSERT INTO members (uuid, name, is_active)
		SELECT * FROM UNNEST($2::uuid[], $3::varchar[], $4::boolean[])
		ON CONFLICT (uuid) DO UPDATE
		SET name = EXCLUDED.name,
		    is_active = EXCLUDED.is_active
		RETURNING id, uuid, name, is_active;
	`

	LinkMembersToTeam = `
		INSERT INTO members_teams (team_id, member_id)
		SELECT $1, m.id
		FROM members m
		WHERE m.uuid = ANY($2::uuid[])
		ON CONFLICT (team_id, member_id) DO NOTHING;
	`

	GetMembersByTeamName = `
		SELECT m.id, m.uuid, m.name, m.is_active
		FROM members m
		INNER JOIN members_teams mt ON m.id = mt.member_id
		INNER JOIN teams t ON mt.team_id = t.id
		WHERE t.name = $1
		ORDER BY m.name;
	`

	GetTeamNameByMemberId = `
		SELECT t.name
		FROM teams t
		INNER JOIN members_teams mt ON t.id = mt.team_id
		INNER JOIN members m ON mt.member_id = m.id
		WHERE m.uuid = $1
		LIMIT 1;
	`
)
