-- name: CreateGroup :one
INSERT INTO contact_groups (id, project_id, name, description)
VALUES (@id, @project_id, @name, @description)
RETURNING *;

-- name: GetGroup :one
SELECT * FROM contact_groups WHERE id = @id AND project_id = @project_id;

-- name: ListGroups :many
SELECT g.*, (SELECT count(*) FROM contact_group_members m WHERE m.group_id = g.id)::bigint AS member_count
FROM contact_groups g
WHERE g.project_id = @project_id
ORDER BY lower(g.name);

-- name: UpdateGroup :one
UPDATE contact_groups SET name = @name, description = @description, updated_at = now()
WHERE id = @id AND project_id = @project_id
RETURNING *;

-- name: DeleteGroup :execrows
DELETE FROM contact_groups WHERE id = @id AND project_id = @project_id;

-- name: AddGroupMembers :execrows
INSERT INTO contact_group_members (group_id, contact_id)
SELECT @group_id, c.id FROM contacts c
WHERE c.project_id = @project_id AND c.id = ANY(@contact_ids::uuid[])
ON CONFLICT DO NOTHING;

-- name: RemoveGroupMember :execrows
DELETE FROM contact_group_members WHERE group_id = @group_id AND contact_id = @contact_id;

-- name: ListGroupMembers :many
SELECT c.* FROM contacts c
JOIN contact_group_members m ON m.contact_id = c.id
WHERE m.group_id = @group_id
ORDER BY lower(c.name), c.created_at
LIMIT @row_limit OFFSET @row_offset;

-- name: CountGroupMembers :one
SELECT count(*)::bigint FROM contact_group_members WHERE group_id = @group_id;

-- name: ListGroupContactIDs :many
SELECT m.contact_id FROM contact_group_members m
JOIN contact_groups g ON g.id = m.group_id
WHERE m.group_id = ANY(@group_ids::uuid[]) AND g.project_id = @project_id;

-- name: ListGroupsForContact :many
SELECT g.* FROM contact_groups g
JOIN contact_group_members m ON m.group_id = g.id
WHERE m.contact_id = @contact_id
ORDER BY lower(g.name);
