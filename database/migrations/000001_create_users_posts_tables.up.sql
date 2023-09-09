CREATE SCHEMA IF NOT EXISTS cms;
CREATE TABLE IF NOT EXISTS cms.post (
    id int PRIMARY KEY,
    title text NOT NULL,
    content text NOT NULL,
    shortdescription text NOT NULL,
    createdat timestamptz NOT NULL,
    urltitle text NOT NULL
);
CREATE TABLE IF NOT EXISTS cms.user (
    id int PRIMARY KEY,
    username text UNIQUE NOT NULL,
    password text NOT NULL,
    email text NOT NULL,
    createdat timestamptz NOT NULL,
    role int NOT NULL
);