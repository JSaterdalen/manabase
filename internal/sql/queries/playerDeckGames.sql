-- name: GetPlayerDeckGame :many
SELECT
    player.id AS player_id,
    player."name" AS player_name,
    deck.id AS deck_id,
    deck."name" AS deck_name,
    pdg.is_won,
    game.id AS game_id,
    game.date_played,
    game.is_totem
FROM
    player_deck_game pdg
    JOIN player ON pdg.player_id = player.id
    JOIN game ON pdg.game_id = game.id
    JOIN deck ON pdg.deck_id = deck.id
ORDER BY
    date_played DESC, game_id;

-- name: InsertPlayerDeckGames :exec
INSERT INTO player_deck_game (id, player_id, game_id, deck_id, is_won, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7);
