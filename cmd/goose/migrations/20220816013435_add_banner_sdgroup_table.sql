-- +goose Up
CREATE TABLE sd_groups
(
    id          serial PRIMARY KEY,
    description text
);

-- +goose Down
DROP TABLE sd_groups;
