CREATE TABLE IF NOT EXISTS likes (
    liked_id integer not null,
    liker_id integer not null,
    CONSTRAINT likes_pkey PRIMARY KEY (liked_id, liker_id),
    CONSTRAINT fk_liked FOREIGN KEY (liked_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_liker FOREIGN KEY (liker_id) REFERENCES users(id) ON DELETE CASCADE
)