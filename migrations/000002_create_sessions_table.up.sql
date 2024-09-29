CREATE TABLE IF NOT EXISTS sessions (
    user_id integer not null,
    uuid    text    not null,
    CONSTRAINT auths_pkey PRIMARY KEY (user_id, uuid)
)