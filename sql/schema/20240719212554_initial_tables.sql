-- +goose Up
-- +goose StatementBegin
CREATE TABLE player (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz,
    name TEXT NOT NULL
);

CREATE TABLE deck (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz,
    name TEXT NOT NULL,
    commander TEXT,
    owner_id UUID NOT NULL REFERENCES player(id)
);

CREATE TABLE game (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz,
    date_played DATE NOT NULL,
    game_number INT NOT NULL,
    is_totem bool NOT NULL DEFAULT FALSE
);

CREATE TABLE player_deck_game (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz,
    player_id UUID NOT NULL REFERENCES player(id),
    game_id UUID NOT NULL REFERENCES game(id),
    deck_id UUID NOT NULL REFERENCES deck(id),
    is_won bool NOT NULL DEFAULT FALSE,
    UNIQUE (player_id, deck_id, game_id)
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE deck,
player,
game,
player_deck_game;

-- +goose StatementEnd