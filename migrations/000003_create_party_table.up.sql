CREATE TABLE IF NOT EXISTS parties (
   id bigserial NOT NULL PRIMARY KEY,
   title text NOT NULL,
   description text NOT NULL,
   image text null,
   price integer DEFAULT 0,
   start_date timestamp,
   creator_id bigint NOT NULL
);