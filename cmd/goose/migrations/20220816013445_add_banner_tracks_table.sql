-- +goose Up
CREATE TABLE tracks
(
    id          serial PRIMARY KEY,
    banner_id   int NOT NULL,
    slot_id     int NOT NULL,
    sd_group_id int NOT NULL,
    clicks      int NOT NULL,
    views       int NOT NULL
);

-- +goose Down
DROP TABLE tracks;