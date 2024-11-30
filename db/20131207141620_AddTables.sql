-- +goose Up

CREATE TABLE IF NOT EXISTS users (
    id          SERIAL PRIMARY KEY,
    group_id    INTEGER DEFAULT 0,
    created_on  TIMESTAMP NOT NULL,
    username    VARCHAR(20),
    password    VARCHAR(75),
    reputations INTEGER DEFAULT 0,
    experience  INTEGER DEFAULT 0,
    level	INTEGER DEFAULT 0,
    avatar      VARCHAR,
    salt        VARCHAR(25)
); CREATE TABLE IF NOT EXISTS boards (
    id           SERIAL PRIMARY KEY,
    title        VARCHAR(45),
    description  VARCHAR(140)
);

CREATE TABLE IF NOT EXISTS posts (
    id           SERIAL PRIMARY KEY,
    board_id     INTEGER REFERENCES boards(id) NOT NULL,
    parent_id    INTEGER REFERENCES posts(id),
    author_id    INTEGER REFERENCES users(id) NOT NULL,
    title        VARCHAR(70) NOT NULL,
    content      TEXT NOT NULL,
    created_on   TIMESTAMP NOT NULL,
    latest_reply TIMESTAMP,
    views        INTEGER DEFAULT 1,
    last_edit    TIMESTAMP,
    sticky       BOOLEAN DEFAULT 'NO'
);

CREATE TABLE IF NOT EXISTS followers (
    id SERIAL PRIMARY KEY,
    follower_id INT NOT NULL,
    followed_id INT NOT NULL,
    FOREIGN KEY (follower_id) REFERENCES users(id),
    FOREIGN KEY (followed_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS notifications (
    id 		  SERIAL PRIMARY KEY,
    user_id 	  INT NOT NULL,
    notif_user_id INT NOT NULL,
    read          BOOLEAN DEFAULT FALSE,
    message 	  TEXT,
    created_on 	  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY   (user_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS notifications (
    id 		  SERIAL PRIMARY KEY,
    user_id 	  INT NOT NULL,
    message_user_id INT NOT NULL,
    read          BOOLEAN DEFAULT FALSE,
    message 	  TEXT,
    created_on 	  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY   (user_id) REFERENCES users(id)
);

-- +goose Down

DROP TABLE users;
DROP TABLE boards;
DROP TABLE posts;
