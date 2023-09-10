CREATE SCHEMA IF NOT EXISTS cms;
CREATE TABLE IF NOT EXISTS cms.post (
    id bigint PRIMARY KEY,
    title text NOT NULL,
    content text NOT NULL,
    shortdescription text NOT NULL,
    createdat timestamptz NOT NULL,
    urltitle text NOT NULL,
    draft boolean DEFAULT true
);
CREATE TABLE IF NOT EXISTS cms.user (
    id bigint PRIMARY KEY,
    username text UNIQUE NOT NULL,
    password text NOT NULL,
    email text NOT NULL,
    createdat timestamptz NOT NULL,
    role text NOT NULL
);
INSERT INTO cms.user(id, username, password, email, createdat, role)
VALUES (
        1337,
        'admin',
        '$2a$10$sXftBHuRTOk.G3JFz2bekOTeLfSen8gpy6O4R2DPtj9PCSpX.HRWC',
        'admin@microblogger.com',
        now(),
        'admin'
    );