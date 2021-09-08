DROP TABLE IF EXISTS public.accounts;
CREATE TABLE IF NOT EXISTS public.accounts
(
    uuid          uuid      NOT NULL PRIMARY KEY default gen_random_uuid(),
    email         text      NOT NULL UNIQUE,
    password_hash text      NOT NULL,
    created_at    timestamp NOT NULL             default current_timestamp,
    updated_at    timestamp,
    deleted_at    timestamp
);

DROP TABLE IF EXISTS public.users;
CREATE TABLE IF NOT EXISTS public.users
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


DROP TABLE IF EXISTS public.posts;
CREATE TABLE IF NOT EXISTS public.posts
(
    uuid       uuid UNIQUE NOT NULL default gen_random_uuid(),
    user_uuid  uuid        NOT NULL,
    title      text        NOT NULL,
    body       text        NOT NULL,
    created_at timestamp   NOT NULL default current_timestamp,
    updated_at timestamp,
    deleted_at timestamp
);
CREATE INDEX ON posts (uuid, created_at);

DROP TABLE IF EXISTS public.comments;
CREATE TABLE IF NOT EXISTS public.comments
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
CREATE INDEX ON comments (created_at DESC);