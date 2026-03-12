-- name: GetPostForUser :many

SELECT posts.* FROM posts 
  INNER JOIN feeds ON feeds.id = posts.feed_id
  INNER JOIN users ON users.id = feeds.user_id
WHERE users.id = $1
ORDER BY posts.updated_at DESC;
