-- name: InsertRepo :exec
INSERT INTO repositories (
        full_name,
        description,
        stargazers_count,
        forks_count,
        created_at
    )
VALUES ($1, $2, $3, $4, $5) ON CONFLICT (full_name) DO
UPDATE
SET description = EXCLUDED.description,
    stargazers_count = EXCLUDED.stargazers_count,
    forks_count = EXCLUDED.forks_count,
    created_at = EXCLUDED.created_at;

-- name: GetRepo :one
SELECT *
FROM repositories
WHERE full_name = $1
LIMIT 1;

-- name: ListAllRepos :many
SELECT *
FROM repositories
ORDER BY full_name;