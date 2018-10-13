#!/bin/bash

createdb -O forum forum
psql -d forum -c "CREATE EXTENSION IF NOT EXISTS citext;"