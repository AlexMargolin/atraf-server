/*ACCOUNTS*/
DROP TABLE IF EXISTS accounts;
CREATE TABLE IF NOT EXISTS accounts
(
    uuid          uuid      NOT NULL PRIMARY KEY default gen_random_uuid(),
    email         text      NOT NULL UNIQUE,
    password_hash text      NOT NULL,
    created_at    timestamp NOT NULL             default current_timestamp,
    updated_at    timestamp,
    deleted_at    timestamp
);

DROP TABLE IF EXISTS accounts_reset;
CREATE TABLE IF NOT EXISTS accounts_reset
(
    account_uuid uuid        NOT NULL PRIMARY KEY,
    token_uuid   uuid UNIQUE NOT NULL default gen_random_uuid(),
    created_at   timestamp   NOT NULL default current_timestamp
);

/*USERS*/
DROP TABLE IF EXISTS users;
CREATE TABLE IF NOT EXISTS users
(
    uuid            uuid      NOT NULL PRIMARY KEY default gen_random_uuid(),
    account_uuid    uuid      NOT NULL UNIQUE,
    email           text UNIQUE,
    first_name      text,
    last_name       text,
    profile_picture text,
    created_at      timestamp NOT NULL             default current_timestamp,
    updated_at      timestamp,
    deleted_at      timestamp
);

/*POSTS*/
DROP TABLE IF EXISTS posts;
CREATE TABLE IF NOT EXISTS posts
(
    uuid       uuid PRIMARY KEY NOT NULL default gen_random_uuid(),
    user_uuid  uuid             NOT NULL,
    title      text             NOT NULL,
    body       text             NOT NULL,
    created_at timestamp        NOT NULL default current_timestamp,
    updated_at timestamp,
    deleted_at timestamp
);
DROP INDEX IF EXISTS posts_created_at_uuid_idx;
CREATE INDEX posts_created_at_uuid_idx ON posts (created_at, uuid);

/*COMMENTS*/
DROP TABLE IF EXISTS comments;
CREATE TABLE IF NOT EXISTS comments
(
    uuid        uuid      NOT NULL PRIMARY KEY default gen_random_uuid(),
    user_uuid   uuid      NOT NULL, /*index*/
    source_uuid uuid      NOT NULL, /*index*/
    parent_uuid uuid      NOT NULL,
    body        text,
    created_at  timestamp NOT NULL             default current_timestamp,
    updated_at  timestamp,
    deleted_at  timestamp
);
DROP INDEX IF EXISTS posts_created_at_idx;
CREATE INDEX posts_created_at_idx ON comments (created_at DESC);