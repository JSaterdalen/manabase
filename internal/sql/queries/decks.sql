-- name: GetDeck :one
SELECT * FROM deck WHERE id = $1;

-- name: GetDeckByName :one
SELECT * FROM deck WHERE name = $1;

-- name: GetDecks :many
SELECT * FROM deck;

-- name: GetDecksByPlayerId :many
SELECT * FROM deck WHERE owner_id = $1;

-- name: GetPlayerDecks :many
SELECT * FROM deck WHERE owner_id = ANY($1::uuid[]);

-- name: CreateDeck :one
INSERT INTO deck (updated_at, name, commander, owner_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetDecksByLastPlayed :many
SELECT
	d.id, 
	d.name,
    d.owner_id,
	MAX(g.date_played) AS last_played_date 
FROM deck d
LEFT JOIN player_deck_game pdg
  ON d.id = pdg.deck_id
LEFT JOIN game g
  ON pdg.game_id = g.id
WHERE d.owner_id = ANY($1::uuid[])
GROUP BY d.id
ORDER BY last_played_date DESC NULLS LAST;

-- name: GetDecksByLastPlayedByPlayer :many
SELECT
	d.id, 
	d.name,
    d.owner_id,
	MAX(g.date_played) AS last_played_date 
FROM deck d
LEFT JOIN player_deck_game pdg
  ON d.id = pdg.deck_id
  AND pdg.player_id = @player::uuid
LEFT JOIN game g
  ON pdg.game_id = g.id
WHERE d.owner_id = ANY(@playerIDs::uuid[])
GROUP BY d.id
ORDER BY last_played_date DESC NULLS LAST;