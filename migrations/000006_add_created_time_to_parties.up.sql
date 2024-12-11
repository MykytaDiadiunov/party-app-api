ALTER TABLE parties
ADD COLUMN created_date timestamp NOT NULL DEFAULT NOW();