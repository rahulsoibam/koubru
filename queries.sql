-- SEPARATE CREDENTIALS DATABASE FOR SECURITY
CREATE DATABASE auth;

CREATE TABLE Credential (
    user_id bigint PRIMARY KEY NOT NULL,
    password_hash text NOT NULL
);

CREATE TABLE Session (
    user_id bigint NOT NULL,
    token text PRIMARY KEY NOT NULL,
    user_agent text NOT NULL DEFAULT 'Unknown',
    created_on timestamptz DEFAULT (now() at time zone 'utc')
);

-- Production Database

CREATE DATABASE koubru_prod;

-- set time zone to utc

SET TIME ZONE 'UTC';

-- global sequence generator

CREATE SEQUENCE global_id_sequence
;

CREATE OR REPLACE FUNCTION id_generator (OUT result bigint)
AS $$
DECLARE
    our_epoch bigint := 1314220021721;
    seq_id bigint;
    now_millis bigint;
    -- the id of this DB shard, must be set for each
    -- schema shard you have - you could pass this as a parameter too
    shard_id int := 1;
BEGIN
    SELECT
        nextval('global_id_sequence') % 1024 INTO seq_id;
    SELECT
        FLOOR(EXTRACT(EPOCH FROM clock_timestamp()) * 1000) INTO now_millis;
    result := (now_millis - our_epoch) << 23;
    result := result | (shard_id << 10);
    result := result | (seq_id);
END;
$$
LANGUAGE PLPGSQL;

SELECT
    id_generator ();

-- reaction types

CREATE TYPE reaction_type AS ENUM ( 'happy',
    'angry',
    'sad',
    'wow',
    'love'
);

-- user

CREATE TABLE KUser (
    user_id bigint PRIMARY KEY DEFAULT id_generator (),
    username citext NOT NULL UNIQUE,
    full_name text NOT NULL,
    bio text DEFAULT '',
    email text UNIQUE DEFAULT NULL,
    phone text UNIQUE DEFAULT NULL,
    email_verified boolean NOT NULL DEFAULT FALSE,
    phone_verified boolean NOT NULL DEFAULT FALSE,
    picture text NOT NULL DEFAULT '_blank',
    facebook_uid text UNIQUE DEFAULT NULL,
    google_uid text UNIQUE DEFAULT NULL,
    created_on timestamptz DEFAULT (now() at time zone 'utc')
);

-- User_user

CREATE TABLE User_Follower (
    user_id bigint NOT NULL REFERENCES KUser (user_id),
    follower_id bigint NOT NULL REFERENCES KUser (user_id),
    followed_on timestamptz DEFAULT (now() at time zone 'utc'),
    PRIMARY KEY (user_id, follower_id)
);

-- Category

CREATE EXTENSION citext;

CREATE TABLE Category (
    category_id bigint PRIMARY KEY DEFAULT id_generator (),
    name citext NOT NULL UNIQUE,
    creator_id bigint NOT NULL REFERENCES KUser (user_id),
    created_on timestamptz DEFAULT (now() at time zone 'utc')
);

-- Category_Follower

CREATE TABLE Category_Follower (
    category_id bigint REFERENCES Category (category_id),
    follower_id bigint REFERENCES KUser (user_id),
    PRIMARY KEY (category_id, user_id),
    followed_on timestamptz DEFAULT (now() at time zone 'utc')
);

-- Topic

CREATE TABLE Topic (
    topic_id bigint PRIMARY KEY DEFAULT id_generator (),
    creator_id bigint NOT NULL REFERENCES KUser (user_id),
    title text NOT NULL,
    details text NOT NULL DEFAULT '',
    created_on timestamptz DEFAULT (now() at time zone 'utc')
);

-- Topic_Follower

CREATE TABLE Topic_Follower (
    topic_id bigint NOT NULL REFERENCES Topic (topic_id),
    follower_id bigint NOT NULL REFERENCES KUser (user_id),
    followed_on timestamptz DEFAULT (now() at time zone 'utc'),
    PRIMARY KEY (topic_id, follower_id)
);

-- Opinion

CREATE TABLE Opinion (
    opinion_id bigint PRIMARY KEY DEFAULT id_generator (),
    parent_id bigint REFERENCES Opinion (opinion_id) DEFAULT NULL,
    creator_id bigint NOT NULL REFERENCES KUser (user_id) NOT NULL,
    topic_id bigint NOT NULL REFERENCES Topic (topic_id) NOT NULL,
    is_anonymous boolean NOT NULL DEFAULT FALSE,
    thumbnails text[] NOT NULL default '{"_blank"}',
    dash text NOT NULL DEFAULT '_blank',
    hls text NOT NULL DEFAULT '_blank',
    created_on timestamptz DEFAULT (now() at time zone 'utc'),
    reaction reaction_type NOT NULL
);

-- Opinion_View

CREATE TABLE Opinion_View (
    opinion_id bigint NOT NULL REFERENCES Opinion (opinion_id),
    viewer_id bigint REFERENCES KUser (user_id),
    viewed_on timestamptz DEFAULT (now() at time zone 'utc')
);

-- Opinion_Vote

CREATE TABLE Opinion_Vote (
    opinion_id bigint NOT NULL REFERENCES Opinion (opinion_id),
    voter_id bigint NOT NULL REFERENCES KUser (user_id),
    vote boolean,
    PRIMARY KEY (opinion_id, voter_id),
    voted_on timestamptz DEFAULT (now() at time zone 'utc')
);

-- Opinion_Follower

CREATE TABLE Opinion_Follower (
    opinion_id bigint REFERENCES Opinion (opinion_id),
    follower_id bigint REFERENCES KUser (user_id),
    PRIMARY KEY (opinion_id, follower_id),
    followed_on timestamptz DEFAULT (now() at time zone 'utc')
);

-- Topic_Category

CREATE TABLE Topic_Category (
    topic_id bigint NOT NULL REFERENCES Topic (topic_id),
    category_id bigint NOT NULL REFERENCES Category (category_id),
    created_on timestamptz DEFAULT (now() at time zone 'utc'),
    PRIMARY KEY (topic_id, category_id)
);

--- Breadcrumbs
WITH RECURSIVE ancestors as (
            SELECT opinion_id, parent_id, creator_id, created_on
            FROM Opinion
            WHERE opinion_id = 1972552364342641740
          UNION ALL
            SELECT o.opinion_id, o.parent_id, o.creator_id, o.created_on
            FROM Opinion AS o
            JOIN ancestors
            ON o.opinion_id = ancestors.parent_id
        )
SELECT
        a.opinion_id,
        u.username,
        u.full_name,
        u.picture,
        case when u.user_id = 1967600534613394434 then 1 else 0 end as is_self,
        (select count(*) from opinion where parent_id=a.opinion_id) as reply_count
FROM ancestors a INNER JOIN Kuser u ON u.user_id = a.creator_id order by a.created_on asc;