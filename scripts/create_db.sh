#!/bin/bash

psql --command "CREATE USER apiforum WITH SUPERUSER PASSWORD 'apiforum';"
createdb -O apiforum apiforum
psql -d apiforum -c "CREATE EXTENSION IF NOT EXISTS citext;"
psql apiforum -f ../sql/init.sql