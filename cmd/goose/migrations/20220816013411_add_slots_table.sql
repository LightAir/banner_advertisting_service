-- +goose Up
CREATE TABLE slots
(
    id          serial PRIMARY KEY,
    description text
);

-- +goose Down
DROP TABLE slots;
