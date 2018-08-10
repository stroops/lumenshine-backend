#!/bin/bash

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" <<-EOSQL
    CREATE USER icop WITH PASSWORD 'jw8s0F4';
    GRANT ALL PRIVILEGES ON DATABASE postgres TO icop;

    CREATE DATABASE icop
      WITH
      OWNER = icop
      ENCODING = 'UTF8'
      LC_COLLATE = 'en_US.utf8'
      LC_CTYPE = 'en_US.utf8'
      TABLESPACE = pg_default
      CONNECTION LIMIT = -1;
    GRANT ALL PRIVILEGES ON DATABASE icop TO icop;    

    CREATE DATABASE admin
      WITH
      OWNER = icop
      ENCODING = 'UTF8'
      LC_COLLATE = 'en_US.utf8'
      LC_CTYPE = 'en_US.utf8'
      TABLESPACE = pg_default
      CONNECTION LIMIT = -1;
    GRANT ALL PRIVILEGES ON DATABASE admin TO icop;

    CREATE DATABASE chart
      WITH
      OWNER = icop
      ENCODING = 'UTF8'
      LC_COLLATE = 'en_US.utf8'
      LC_CTYPE = 'en_US.utf8'
      TABLESPACE = pg_default
      CONNECTION LIMIT = -1;
    GRANT ALL PRIVILEGES ON DATABASE chart TO icop;

    CREATE DATABASE dividend
      WITH
      OWNER = icop
      ENCODING = 'UTF8'
      LC_COLLATE = 'en_US.utf8'
      LC_CTYPE = 'en_US.utf8'
      TABLESPACE = pg_default
      CONNECTION LIMIT = -1;
    GRANT ALL PRIVILEGES ON DATABASE dividend TO icop;
EOSQL

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" -d icop <<-EOSQL
    CREATE EXTENSION pg_trgm;
EOSQL

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" -d admin <<-EOSQL
    CREATE EXTENSION pg_trgm;
EOSQL

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" -d chart <<-EOSQL
    CREATE EXTENSION pg_trgm;
EOSQL

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" -d dividend <<-EOSQL
    CREATE EXTENSION pg_trgm;
EOSQL