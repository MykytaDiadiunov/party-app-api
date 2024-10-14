CREATE TABLE IF NOT EXISTS party_users (
    party_id integer not null,
    user_id integer not null,
    CONSTRAINT party_users_pkey PRIMARY KEY (party_id, user_id)
)