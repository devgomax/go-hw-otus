-- +goose Up
-- +goose StatementBegin
ALTER TABLE events
    ADD COLUMN processed BOOLEAN DEFAULT FALSE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE events
    DROP COLUMN IF EXISTS processed;
-- +goose StatementEnd
