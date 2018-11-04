FROM ubuntu:16.04
MAINTAINER Daniel Lee

# PostgreSQL installing
RUN apt-get -y update
RUN apt-get -y install apt-transport-https git wget
RUN echo 'deb http://apt.postgresql.org/pub/repos/apt/ xenial-pgdg main' >> /etc/apt/sources.list.d/pgdg.list
RUN wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add -
RUN apt-get -y update
ENV PGVERSION 10
RUN apt-get -y install postgresql-$PGVERSION postgresql-contrib

# Repo
USER root
RUN git clone https://github.com/Unanoc/tpark_db.git
WORKDIR tpark_db

# PostgreSQL creating of database
USER postgres
RUN /etc/init.d/postgresql start &&\
    psql --command "CREATE USER apiforum WITH SUPERUSER PASSWORD 'apiforum';" &&\
    createdb -O apiforum apiforum &&\
    psql -d apiforum -c "CREATE EXTENSION IF NOT EXISTS citext;" &&\
    psql apiforum -a -f sql/create_tables.sql &&\
    /etc/init.d/postgresql stop
    
USER root
# Open Postgres for network
RUN echo "local all all md5" > /etc/postgresql/$PGVERSION/main/pg_hba.conf &&\
    echo "host all all 0.0.0.0/0 md5" >> /etc/postgresql/$PGVERSION/main/pg_hba.conf &&\
    echo "listen_addresses='*'" >> /etc/postgresql/$PGVERSION/main/postgresql.conf &&\
    echo "unix_socket_directories = '/var/run/postgresql'" >> /etc/postgresql/$PGVERSION/main/postgresql.conf
VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]
EXPOSE 5432

# Golang installing
ENV GOVERSION 1.11.1
USER root
RUN wget https://storage.googleapis.com/golang/go$GOVERSION.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go$GOVERSION.linux-amd64.tar.gz && \
    mkdir go && mkdir go/src && mkdir go/bin && mkdir go/pkg
ENV GOROOT /usr/local/go
ENV GOPATH /opt/go
ENV PATH $GOROOT/bin:$GOPATH/bin:$PATH
RUN mkdir -p "$GOPATH/bin" "$GOPATH/src"
RUN apt-get -y install gcc musl-dev && GO11MODULE=on
RUN go build .
EXPOSE 5000
RUN echo "./config/postgresql.conf" >> /etc/postgresql/$PGVERSION/main/postgresql.conf

USER root
CMD service postgresql start && go run main.go
