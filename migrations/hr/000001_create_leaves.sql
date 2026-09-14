-- +goose Up
CREATE TABLE leaves (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    start_date TEXT NOT NULL,
    end_date TEXT NOT NULL
);

-- +goose Down
DROP TABLE leaves;
