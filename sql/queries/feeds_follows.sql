-- name: CreateFeedFollow :one
WITH inserted_feed_follows AS (
    INSERT INTO feed_follows(user_id, feed_id)
    VALUES ($1, $2)
    ON CONFLICT (user_id, feed_id) DO NOTHING
    RETURNING *
),
picked AS(
    SELECT * FROM inserted_feed_follows
    UNION ALL
    SELECT feed_follows.*
    FROM feed_follows
    WHERE NOT EXISTS (SELECT 1 FROM inserted_feed_follows)
        AND feed_follows.user_id = $1 AND feed_follows.feed_id = $2
)
SELECT
    picked.user_id AS follow_user_id,
    picked.feed_id AS follow_feed_id,
    feeds.name AS feed_name,
    users.name AS user_name
FROM picked
INNER JOIN users ON picked.user_id = users.id
INNER JOIN feeds ON picked.feed_id = feeds.id;

-- name: GetAllFeedFollow :many
SELECT * FROM feed_follows;

-- name: GetFeedFollowForUser :many
SELECT
    feed_follows.user_id AS  follow_user_id,
    feed_follows.feed_id AS  follow_feed_id,
    users.name AS user_name,
    feeds.name AS feed_name
FROM feed_follows
INNER JOIN users ON feed_follows.user_id = users.id
INNER JOIN feeds ON feed_follows.feed_id = feeds.id
WHERE feed_follows.user_id = $1;

-- name: DeletFeedFollowByURL :exec
DELETE FROM feed_follows WHERE feed_follows.user_id = $1 AND feed_follows.feed_id = $2;

