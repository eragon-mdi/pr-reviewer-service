package queries

const (
	CreatePullRequest = `
		WITH pr_ins AS (
			INSERT INTO pull_requests (uuid, title, author_id, status_id, created_at, version)
			VALUES ($1, $2, $3, (SELECT id FROM statuses WHERE status = 'OPEN'), NOW(), 1)
			ON CONFLICT (uuid) DO NOTHING
			RETURNING id, author_id
		),
		pr_existing AS (
			SELECT id, author_id
			FROM pull_requests
			WHERE uuid = $1
		),
		pr_id AS (
			SELECT id, author_id FROM pr_ins
			UNION ALL
			SELECT id, author_id FROM pr_existing
			WHERE NOT EXISTS (SELECT 1 FROM pr_ins)
		),
		reviewers AS (
			SELECT m.id, m.uuid
			FROM members m
			INNER JOIN members_teams mt ON m.id = mt.member_id
			INNER JOIN pr_id ON mt.team_id = (
				SELECT mt2.team_id 
				FROM members_teams mt2
				INNER JOIN members m2 ON mt2.member_id = m2.id
				WHERE m2.id = (SELECT author_id FROM pr_id LIMIT 1)
				LIMIT 1
			)
			WHERE m.is_active = true
			  AND m.id != (SELECT author_id FROM pr_id LIMIT 1)
			ORDER BY RANDOM()
			LIMIT 2
		)
		INSERT INTO pr_members (pr_id, member_id, role_id, assigned_at)
		SELECT 
			(SELECT id FROM pr_id LIMIT 1),
			r.id,
			(SELECT id FROM roles WHERE role = 'reviewer'),
			NOW()
		FROM reviewers r
		ON CONFLICT (pr_id, member_id) DO NOTHING
		RETURNING pr_id, member_id;
	`

	GetPullRequestByUUID = `
		SELECT 
			pr.id,
			pr.uuid,
			pr.title,
			author.uuid AS author_id,
			s.status AS status,
			pr.created_at,
			pr.merged_at,
			pr.version
		FROM pull_requests pr
		INNER JOIN members author ON pr.author_id = author.id
		INNER JOIN statuses s ON pr.status_id = s.id
		WHERE pr.uuid = $1;
	`

	GetPullRequestReviewers = `
		SELECT m.id, m.uuid, m.name, m.is_active
		FROM pr_members pm
		INNER JOIN pull_requests pr ON pm.pr_id = pr.id
		INNER JOIN members m ON pm.member_id = m.id
		INNER JOIN roles r ON pm.role_id = r.id
		WHERE pr.uuid = $1
		  AND r.role = 'reviewer'
		ORDER BY pm.assigned_at;
	`

	MergePullRequest = `
		UPDATE pull_requests
		SET status_id = (SELECT id FROM statuses WHERE status = 'MERGED'),
		    merged_at = COALESCE(merged_at, NOW()),
		    version = version + 1
		WHERE uuid = $1
		  AND status_id != (SELECT id FROM statuses WHERE status = 'MERGED')
		RETURNING id, uuid, title, author_id, status_id, created_at, merged_at, version;
	`

	GetPullRequestMembersHistories = `
		SELECT 
			m.id,
			m.uuid,
			m.name,
			m.is_active,
			r.role,
			pm.assigned_at,
			CASE 
				WHEN pm.member_id = $2 THEN true 
				ELSE false 
			END AS was_assigned_before
		FROM members m
		INNER JOIN members_teams mt ON m.id = mt.member_id
		INNER JOIN pull_requests pr ON mt.team_id = (
			SELECT mt2.team_id 
			FROM members_teams mt2
			INNER JOIN members m2 ON mt2.member_id = m2.id
			WHERE m2.id = (
				SELECT pm2.member_id 
				FROM pr_members pm2
				INNER JOIN pull_requests pr2 ON pm2.pr_id = pr2.id
				WHERE pr2.uuid = $1
				LIMIT 1
			)
			LIMIT 1
		)
		LEFT JOIN pr_members pm ON pm.member_id = m.id AND pm.pr_id = (
			SELECT id FROM pull_requests WHERE uuid = $1
		)
		LEFT JOIN roles r ON pm.role_id = r.id
		WHERE EXISTS (
			SELECT 1 FROM pull_requests WHERE uuid = $1
		)
		ORDER BY m.name;
	`

	AssignMemberToPR = `
		WITH old_reviewer AS (
			DELETE FROM pr_members
			WHERE pr_id = (SELECT id FROM pull_requests WHERE uuid = $1)
			  AND member_id = (SELECT id FROM members WHERE uuid = $2)
			  AND role_id = (SELECT id FROM roles WHERE role = 'reviewer')
			RETURNING pr_id
		),
		new_assignment AS (
			INSERT INTO pr_members (pr_id, member_id, role_id, assigned_at)
			SELECT 
				(SELECT id FROM pull_requests WHERE uuid = $1),
				(SELECT id FROM members WHERE uuid = $3),
				(SELECT id FROM roles WHERE role = 'reviewer'),
				NOW()
			WHERE EXISTS (SELECT 1 FROM old_reviewer)
			RETURNING pr_id
		)
		SELECT 
			pr.id,
			pr.uuid,
			pr.title,
			author.uuid AS author_id,
			s.status AS status,
			pr.created_at,
			pr.merged_at,
			pr.version
		FROM pull_requests pr
		INNER JOIN members author ON pr.author_id = author.id
		INNER JOIN statuses s ON pr.status_id = s.id
		WHERE pr.uuid = $1;
	`

	CheckPRStatus = `
		SELECT s.status
		FROM pull_requests pr
		INNER JOIN statuses s ON pr.status_id = s.id
		WHERE pr.uuid = $1;
	`

	CheckMemberAssignedToPR = `
		SELECT EXISTS(
			SELECT 1
			FROM pr_members pm
			INNER JOIN pull_requests pr ON pm.pr_id = pr.id
			INNER JOIN members m ON pm.member_id = m.id
			INNER JOIN roles r ON pm.role_id = r.id
			WHERE pr.uuid = $1
			  AND m.uuid = $2
			  AND r.role = 'reviewer'
		);
	`
)
