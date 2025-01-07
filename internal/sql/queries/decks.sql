-- name: GetDeck :one
SELECT * FROM deck WHERE id = $1;

-- name: GetDeckByName :one
SELECT * FROM deck WHERE name = $1;

-- name: GetDecks :many
SELECT * FROM deck;

-- name: GetDecksByPlayerId :many
SELECT * FROM deck WHERE owner_id = $1;

-- name: CreateDeck :one
INSERT INTO deck (updated_at, name, commander, owner_id)
VALUES ($1, $2, $3, $4)
RETURNING *;