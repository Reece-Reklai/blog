-- name: CreatePost :one
INSERT INTO posts(title, url, description, published_at, feed_id)
VALUES($1, $2, $3, $4, $5)
ON CONFLICT (url) DO NOTHING
RETURNING *;

-- name: GetPosts :many
SELECT
    posts.title as post_title,
    posts.url as post_url
FROM posts;

-- name: GetPostForUser :many
SELECT
    posts.title as post_title,
    posts.url as post_url,
    posts.description as post_description,
    feeds.name as feed_name,
    feeds.url as feed_url,
    users.name as user_name
FROM posts
INNER JOIN feeds ON feeds.id = posts.feed_id
INNER JOIN feed_follows ON feed_follows.feed_id = feeds.id
INNER JOIN users ON users.id = feed_follows.user_id
WHERE users.id = $1
ORDER BY posts.published_at DESC
LIMIT $2;
