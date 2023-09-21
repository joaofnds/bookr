-- +goose Up
CREATE TABLE resources(
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    own_id varchar,
    setup int8 NOT NULL DEFAULT 0,
    cleanup int8 NOT NULL DEFAULT 0,
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE resources;
