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


DROP TABLE IF EXISTS public.posts;
CREATE TABLE IF NOT EXISTS public.posts
(
    uuid       uuid      NOT NULL PRIMARY KEY default gen_random_uuid(),
    user_uuid  uuid      NOT NULL,
    content    text      NOT NULL,
    created_at timestamp NOT NULL             default current_timestamp,
    updated_at timestamp,
    deleted_at timestamp
);

DROP TABLE IF EXISTS public.comments;
CREATE TABLE IF NOT EXISTS public.comments
(
    uuid        uuid      NOT NULL PRIMARY KEY default gen_random_uuid(),
    user_uuid   uuid      NOT NULL, /*index*/
    post_uuid   uuid      NOT NULL, /*index*/
    parent_uuid uuid      NOT NULL,
    content     text,
    created_at  timestamp NOT NULL             default current_timestamp,
    updated_at  timestamp,
    deleted_at  timestamp
);