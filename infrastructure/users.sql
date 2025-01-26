CREATE TABLE IF NOT EXISTS users
(
    -- https://www.naiyerasif.com/post/2024/09/04/stop-using-serial-in-postgres/
    --id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    --id SERIAL PRIMARY KEY,
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    name TEXT NOT NULL,
    token TEXT NOT NULL
);
