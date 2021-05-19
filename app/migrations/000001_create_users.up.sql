BEGIN;
CREATE TABLE IF NOT EXISTS users(
   id serial PRIMARY KEY,
   username VARCHAR (50) UNIQUE NOT NULL,
   age INT NOT NULL
);
CREATE INDEX age_users on users (age);
COMMIT;