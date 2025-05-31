-- +goose Up
-- +goose StatementBegin
CREATE TABLE events
(
    id              UUID PRIMARY KEY,
    title           TEXT                     NOT NULL,
    starts_at       TIMESTAMP WITH TIME ZONE NOT NULL,
    ends_at         TIMESTAMP WITH TIME ZONE NOT NULL,
    description     TEXT                     NULL,
    user_id         UUID                     NOT NULL,
    notify_interval INTERVAL                 NOT NULL DEFAULT '15 minutes'
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE events;
-- +goose StatementEnd
