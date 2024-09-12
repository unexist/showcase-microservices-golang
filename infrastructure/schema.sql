CREATE TABLE todos
(
    id SERIAL,
    uuid TEXT NOT NULL,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    CONSTRAINT todos_pkey PRIMARY KEY (id)
)
