-- This script is to be ran on local db when initiating a new database
-- On remote database, this is to be ran manually from a superuser (admin), using a different 
-- password for each environment.

----------------------
-- Calidum Rotae --
----------------------

-- Create DB
CREATE DATABASE calidum_rotae
WITH
	LOCALE = 'en_US.UTF-8'
	ENCODING = 'UTF8'
	TEMPLATE = template0;

\connect calidum_rotae

REVOKE CREATE ON SCHEMA public FROM PUBLIC;
REVOKE ALL ON DATABASE calidum_rotae FROM PUBLIC;

-- Create User Schema
CREATE SCHEMA calidum_rotae;

-- Create Users

-- Migration role needs access to whole schema
-- Password only used in local
CREATE ROLE sqlmigrator
WITH
	LOGIN PASSWORD 'some_hard_password'
	NOSUPERUSER;

-- So that cedille_user can do everything sqlmigrator can
-- On Cloud SQL, cedille_user user is not a superuser
GRANT sqlmigrator TO cedille_user;

-- Grant on currently existing tables, if any
GRANT USAGE ON SCHEMA calidum_rotae TO sqlmigrator;
GRANT ALL PRIVILEGES ON SCHEMA  calidum_rotae TO sqlmigrator;
GRANT ALL PRIVILEGES ON DATABASE calidum_rotae TO sqlmigrator;

ALTER DEFAULT PRIVILEGES IN SCHEMA calidum_rotae
GRANT ALL ON TABLES TO sqlmigrator;

GRANT CONNECT ON DATABASE calidum_rotae TO sqlmigrator;
GRANT ALL ON schema public TO sqlmigrator;

-- calidum_rotae is the role used by the calidum rotae service component in the application
-- Password only used in local
CREATE ROLE calidum_rotae
WITH
	LOGIN PASSWORD 'some_password'
	NOSUPERUSER;

-- So that cedille_user can do everything calidum_rotae can
GRANT calidum_rotae TO cedille_user;

GRANT CONNECT ON DATABASE calidum_rotae TO calidum_rotae;

GRANT USAGE ON SCHEMA calidum_rotae TO calidum_rotae;

-- Grant on currently existing tables, if any
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA calidum_rotae TO calidum_rotae;
GRANT USAGE ON ALL SEQUENCES IN SCHEMA calidum_rotae TO calidum_rotae; --is this one required?

-- To grant default privileges on future tables, you need to grant to the user/role you are creating the table with.
ALTER DEFAULT PRIVILEGES FOR ROLE sqlmigrator IN SCHEMA calidum_rotae
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO calidum_rotae;

ALTER DEFAULT PRIVILEGES FOR ROLE sqlmigrator IN SCHEMA calidum_rotae
GRANT USAGE ON SEQUENCES TO calidum_rotae;