-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS products(
  id BIGINT GENERATED BY DEFAULT AS IDENTITY,
  match VARCHAR NOT NULL,
  reward int NOT NULL,
  reward_type VARCHAR NOT NULL,
  PRIMARY KEY(id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
