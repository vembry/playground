CREATE DATABASE poster;

\c poster;

CREATE TABLE public.posts(
    id CHAR(36) PRIMARY KEY,
    user_id VARCHAR(50),
    parent_post_id CHAR(36) NULL,
    content TEXT,
    like_count INT,
    dislike_count INT,
    created_at TIMESTAMP,

    FOREIGN KEY (parent_post_id) REFERENCES posts (id)
)