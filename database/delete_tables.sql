-- revoke ALL on ALL TABLES in SCHEMA if exists f from adm;
\c filmoteka;
drop schema if exists f cascade;
\c postgres;

DROP DATABASE IF EXISTS filmoteka;

drop user if exists adm;
drop user if exists u;
