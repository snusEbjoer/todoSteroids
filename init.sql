CREATE USER todo WITH PASSWORD 'todo123';
CREATE DATABASE tododb OWNER "todo";
GRANT ALL PRIVILEGES ON DATABASE tododb TO "todo";

CREATE USER "user" WITH PASSWORD 'user123';
CREATE DATABASE userdb OWNER "user";
GRANT ALL PRIVILEGES ON DATABASE "userdb" TO "user";

CREATE USER stats WITH PASSWORD 'stats123';
CREATE DATABASE statsdb OWNER "stats";
GRANT ALL PRIVILEGES ON DATABASE statsdb TO "stats";

SELECT usename AS role_name,
 CASE
  WHEN usesuper AND usecreatedb THEN
    CAST('superuser, create database' AS pg_catalog.text)
  WHEN usesuper THEN
    CAST('superuser' AS pg_catalog.text)
  WHEN usecreatedb THEN
    CAST('create database' AS pg_catalog.text)
  ELSE
    CAST('' AS pg_catalog.text)
 END role_attributes
FROM pg_catalog.pg_user
ORDER BY role_name desc;