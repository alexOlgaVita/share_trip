-- +goose Up
-- +goose StatementBegin
CREATE TABLE items (
                       id   UUID PRIMARY KEY,
                       name TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE items;
-- +goose StatementEnd