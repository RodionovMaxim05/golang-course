-- name: CreateSubscription :one
INSERT INTO subscriptions (owner, repo)
VALUES ($1, $2)
RETURNING id, owner, repo, created_at;

-- name: DeleteSubscription :exec
DELETE FROM subscriptions
WHERE owner = $1
    AND repo = $2;

-- name: ListSubscriptions :many
SELECT *
FROM subscriptions
ORDER BY created_at DESC;

-- name: CheckSubscriptionExists :one
SELECT EXISTS (
        SELECT 1
        FROM subscriptions
        WHERE owner = $1
            AND repo = $2
    );