-- name: InsertGame :exec
INSERT INTO game (
    id,
    created_at,
    updated_at,
    date_played
) VALUES (
    @id,
    @created_at,
    @updated_at,
    @date_played
);