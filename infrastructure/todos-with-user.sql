CREATE TABLE IF NOT EXISTS todos
(
    -- https://www.naiyerasif.com/post/2024/09/04/stop-using-serial-in-postgres/
    --id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    id SERIAL PRIMARY KEY,
    user_id INTEGER,
    title TEXT NOT NULL,
    description TEXT NOT NULL
);
