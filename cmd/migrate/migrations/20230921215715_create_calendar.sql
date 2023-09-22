-- +goose Up
CREATE TABLE calendars(
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    resource_id uuid NOT NULL REFERENCES resources(id) ON DELETE CASCADE,
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE calendars;