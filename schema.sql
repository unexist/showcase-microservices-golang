CREATE TABLE todos
(
    id SERIAL,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    CONSTRAINT todos_pkey PRIMARY KEY (id)
)
