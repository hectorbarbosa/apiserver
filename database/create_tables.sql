CREATE USER adm WITH PASSWORD 'pass';
CREATE USER u WITH PASSWORD 'pass';

CREATE DATABASE filmoteka
  WITH OWNER = postgres
  ENCODING = 'UTF8';

\c filmoteka;

create schema f;

CREATE TABLE IF NOT EXISTS f.films (
	film_id SERIAL PRIMARY KEY,
	film_name varchar(150) NOT NULL,
	description varchar(1000),
	release_date char(10) NOT NULL,
    rating smallint
);

CREATE TABLE IF NOT EXISTS f.actors (
	actor_id SERIAL PRIMARY KEY,
	actor_name varchar(255) NOT NULL,
	gender char(1) NOT NULL,
	birth_date char(10) NOT NULL
);

CREATE TABLE IF NOT EXISTS f.roles (
    role_id SERIAL PRIMARY KEY,
	actor_id int REFERENCES f.actors(actor_id),
	film_id int REFERENCES f.films(film_id)
);

CREATE UNIQUE INDEX if not exists flms ON f.films(film_name, release_date);
CREATE UNIQUE INDEX if not exists actrs ON f.actors(actor_name, birth_date);
CREATE UNIQUE INDEX if not exists r ON f.roles(film_id, actor_id);

GRANT CONNECT ON DATABASE filmoteka TO adm;
GRANT CONNECT ON DATABASE filmoteka TO u;
GRANT USAGE ON schema f TO adm;
GRANT USAGE ON schema f TO u;
-- GRANT ALL PRIVILEGES ON schema f TO adm;
-- GRANT INSERT ON ALL TABLES IN schema f TO adm;
GRANT ALL PRIVILEGES ON ALL SEQUENCES in schema f TO adm;
GRANT ALL PRIVILEGES ON ALL TABLES in schema f TO adm;
GRANT SELECT ON ALL TABLES in schema f TO u;


\c postgres;
