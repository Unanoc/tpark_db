#!/bin/bash

psql --command "DROP DATABASE IF EXISTS apiforum;"
psql --command "DROP USER IF EXISTS apiforum;"