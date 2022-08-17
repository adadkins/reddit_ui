CREATE TABLE comments (
    id varchar(10) NOT NULL UNIQUE PRIMARY KEY,
    parent_id varchar(10) NOT NULL,
    top_level_id varchar(10) NOT NULL,
    post_id varchar(10) NOT NULL,
    body varchar(100) NOT NULL,
    author varchar(50) NOT NULL,
    created_at integer NOT NULL,
    updated_at integer NOT NULL
);
