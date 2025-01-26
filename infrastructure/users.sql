CREATE TABLE users
(
    -- https://www.naiyerasif.com/post/2024/09/04/stop-using-serial-in-postgres/
    --id SERIAL PRIMARY KEY,
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name TEXT NOT NULL,
    token TEXT NOT NULL
);
