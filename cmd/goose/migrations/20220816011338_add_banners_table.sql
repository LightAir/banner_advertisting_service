-- +goose Up
CREATE TABLE banners
(
    id          serial PRIMARY KEY,
    description text
);

-- +goose Down
DROP TABLE banners;
