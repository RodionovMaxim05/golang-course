-- name: InsertRepo :exec
INSERT INTO repositories (
        full_name,
        description,
        stargazers_count,
        forks_count,
        created_at,
        repo_status,
        error_code,
        updated_at
    )
VALUES ($1, $2, $3, $4, $5, $6, $7, NOW()) ON CONFLICT (full_name) DO
UPDATE
SET description = EXCLUDED.description,
    stargazers_count = EXCLUDED.stargazers_count,
    forks_count = EXCLUDED.forks_count,
    created_at = EXCLUDED.created_at,
    repo_status = EXCLUDED.repo_status,
    error_code = EXCLUDED.error_code,
    updated_at = NOW();

-- name: UpdateRepoStatus :exec
INSERT INTO repositories (full_name, repo_status, error_code, updated_at)
VALUES ($1, $2, $3, NOW()) ON CONFLICT (full_name) DO
UPDATE
SET repo_status = EXCLUDED.repo_status,
    error_code = EXCLUDED.error_code,
    updated_at = NOW();

-- name: GetRepo :one
SELECT *
FROM repositories
WHERE full_name = $1
LIMIT 1;

-- name: GetReposByNames :many
SELECT *
FROM repositories
WHERE full_name = ANY($1::text []);

-- name: ListAllRepos :many
SELECT *
FROM repositories
ORDER BY full_name;