CREATE TABLE comments (
    id VARCHAR(10) NOT NULL UNIQUE PRIMARY KEY,
    parent_id VARCHAR(10) NOT NULL,
    body VARCHAR(5000) NOT NULL,
    author VARCHAR(50) NOT NULL,
    updated_at INTEGER NOT NULL
);

CREATE TABLE submission(
     id VARCHAR(10) NOT NULL UNIQUE PRIMARY KEY,
     body VARCHAR(5000) NOT NULL,
     author VARCHAR(50) NOT NULL,
     updated_at INTEGER NOT NULL
)
