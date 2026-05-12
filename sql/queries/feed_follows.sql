-- name: CreateFeedFollow :one
WITH inserted_feed_follows AS (
INSERT INTO feed_follows (id, created_at, updated_at, users_id, feed_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *)
SELECT inserted_feed_follows.*, users.name AS user_name, feeds.name AS feed_name
FROM inserted_feed_follows
INNER JOIN users ON inserted_feed_follows.users_id = users.id
INNER JOIN feeds ON inserted_feed_follows.feed_id = feeds.id;

-- name: GetFeedFollowsForUser :many
SELECT feed_follows.*, users.name AS user_name, feeds.name AS feed_name
FROM feed_follows
INNER JOIN users ON feed_follows.users_id = users.id
INNER JOIN feeds ON feed_follows.feed_id = feeds.id
WHERE feed_follows.users_id = $1;