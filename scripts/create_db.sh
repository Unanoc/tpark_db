#!/bin/bash

createdb -O apiforum apiforum
psql -d apiforum -c "CREATE EXTENSION IF NOT EXISTS citext;"