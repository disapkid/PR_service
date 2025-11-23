-- +goose Up
-- +goose StatementBegin
CREATE TABLE teams(
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE
);

CREATE TABLE users(
    id TEXT PRIMARY KEY,
    name TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    team_id BIGINT,
    FOREIGN KEY(team_id) REFERENCES teams(id)
);

CREATE TABLE pull_requests(
    id TEXT PRIMARY KEY,
    name TEXT,
    author_id TEXT,
    status TEXT DEFAULT 'OPEN',
    FOREIGN KEY(author_id) REFERENCES users(id)
);

CREATE TABLE reviewers(
    pr_id TEXT,
    user_id TEXT,
    PRIMARY KEY (pr_id, user_id),
    FOREIGN KEY(pr_id) REFERENCES pull_requests(id),
    FOREIGN KEY(user_id) REFERENCES users(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE reviewers;
DROP TABLE pull_requests;
DROP TABLE users;
DROP TABLE teams;
-- +goose StatementEnd
