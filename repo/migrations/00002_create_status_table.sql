-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS status (
    name VARCHAR(50),
    PRIMARY KEY (name)
);

INSERT INTO status (name) VALUES
('todo'),
('doing'),
('done')
ON CONFLICT (name) DO NOTHING; -- for not essuing errors if the table already exists
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE status;
-- +goose StatementEnd
