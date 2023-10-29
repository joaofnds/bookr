-- +goose Up
CREATE TABLE calendar_events(
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    calendar_id uuid NOT NULL REFERENCES calendars(id) ON DELETE CASCADE,
    name varchar(255) NOT NULL,
    description text NOT NULL,
    status int NOT NULL,
    starts_at timestamp NOT NULL,
    ends_at timestamp NOT NULL,
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE calendar_events;