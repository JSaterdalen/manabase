-- name: GetPlayerDeckGame :many
SELECT
    player.id AS player_id,
    player."name" AS player_name,
    deck.id AS deck_id,
    deck."name" AS deck_name,
    pdg.is_won,
    game.id AS game_id,
    game.date_played,
    game.game_number,
    game.is_totem
FROM
    player_deck_game pdg
    JOIN player ON pdg.player_id = player.id
    JOIN game ON pdg.game_id = game.id
    JOIN deck ON pdg.deck_id = deck.id
ORDER BY
    date_played DESC;
