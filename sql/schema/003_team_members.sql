-- +goose Up
CREATE TABLE team_members (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    team_id UUID NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    UNIQUE(user_id, team_id)
);

-- +goose Down
DROP TABLE team_members;