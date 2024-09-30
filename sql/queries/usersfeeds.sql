-- name: CreateFeedFollow :one
WITH inserted_feed_follows AS ( 
INSERT INTO feed_follows (id, user_id, feed_id)
VALUES (
    $1,
    $2,
    $3
    )
RETURNING *
)

SELECT inserted_feed_follows.*,
    feeds.name AS feed_name,
    users.name AS user_name

FROM inserted_feed_follows
INNER JOIN feeds
ON inserted_feed_follows.feed_id = feeds.id 
INNER JOIN users
ON inserted_feed_follows.user_id = users.id;

-- name: GetFeedFollowsForUser :many
SELECT feed_follows.*,
    feeds.name AS feed_name,
    users.name AS user_name
FROM feed_follows
INNER JOIN feeds ON feed_follows.feed_id = feeds.id
INNER JOIN users ON feed_follows.user_id = users.id
WHERE feed_follows.user_id = $1;

-- name: GetFeedByIDs :one
SELECT * FROM feed_follows
WHERE user_id = $1 and feed_id = $2;

-- name: DeleteFeedFollows :exec
DELETE FROM feed_follows
WHERE user_id = $1 and feed_id = $2;

