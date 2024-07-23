-- name: GetDeck :one
SELECT * FROM deck WHERE id = $1;

-- name: GetDecks :many
SELECT * FROM deck;

-- name: CreateDeck :one
INSERT INTO deck (updated_at, name, commander, owner_id)
VALUES ($1, $2, $3, $4)
RETURNING *;
