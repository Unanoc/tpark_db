#!/bin/bash
psql --command "CREATE USER forum WITH SUPERUSER PASSWORD 'forum';"
createdb -O forum forum
psql -c "GRANT ALL ON DATABASE postgres TO forum;"
psql -d postgres -c "CREATE EXTENSION IF NOT EXISTS citext;"

# psql forum