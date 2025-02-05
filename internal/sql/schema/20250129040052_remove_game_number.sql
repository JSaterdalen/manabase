-- +goose Up
-- +goose StatementBegin
ALTER TABLE game DROP COLUMN game_number;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE game ADD COLUMN game_number INT NOT NULL DEFAULT 1;
-- +goose StatementEnd
