FROM ubuntu:18.04
LABEL author='Daniel Lee'

# Basic tools installing
RUN apt update
RUN apt install -y git vim wget curl

# PostgreSQL installing
RUN apt install -y postgresql-$PGVERSION postgresql-contrib

# Database creating
USER postgres
RUN /etc/init.d/postgresql start &&\
    psql --command "CREATE USER forum WITH SUPERUSER PASSWORD 'forum';" &&\
    createdb -O forum forum &&\
    psql -d forum -c "CREATE EXTENSION IF NOT EXISTS citext;" &&\
    /etc/init.d/postgresql stop


# для использования go mod нужно скачать след утилиты
# RUN apk add --update git gcc musl-dev && GO11MODULE=on go build