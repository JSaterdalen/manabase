-- name: GetPlayer :one
SELECT * FROM player WHERE id = $1;

-- name: GetPlayers :many
SELECT * FROM player;

-- name: CreatePlayer :one
INSERT INTO player (updated_at, name)
VALUES ($1, $2)
RETURNING *;
