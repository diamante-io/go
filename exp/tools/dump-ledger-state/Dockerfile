FROM ubuntu:20.04

ENV diamcircle_CORE_VERSION=18.1.0-760.rc1.dc5f5a392.focal
ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates curl wget gnupg apt-utils
RUN wget -qO - https://apt.diamcircle.org/SDF.asc | APT_KEY_DONT_WARN_ON_DANGEROUS_USAGE=true apt-key add -
RUN echo "deb https://apt.diamcircle.org focal stable" >/etc/apt/sources.list.d/SDF.list
RUN echo "deb https://apt.diamcircle.org focal unstable" >/etc/apt/sources.list.d/SDF-unstable.list
RUN apt-get update -y

RUN apt-get install -y diamcircle-core=${diamcircle_CORE_VERSION} jq
RUN apt-get clean
RUN echo "deb http://apt.postgresql.org/pub/repos/apt/ $(env -i bash -c '. /etc/os-release; echo $VERSION_CODENAME')-pgdg main" | tee /etc/apt/sources.list.d/pgdg.list && \
    wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add - && \
    apt-get update && \
    DEBIAN_FRONTEND="noninteractive" apt-get install -y postgresql-9.6 postgresql-contrib-9.6 postgresql-client-9.6

# Create a PostgreSQL role named `circleci` and then create a database `core` owned by the `circleci` role.
RUN  su - postgres -c "/etc/init.d/postgresql start && psql --command \"CREATE USER circleci WITH SUPERUSER;\" && createdb -O circleci core"

# Adjust PostgreSQL configuration so that remote connections to the
# database are possible.
RUN echo "host all all all trust" > /etc/postgresql/9.6/main/pg_hba.conf

# And add `listen_addresses` to `/etc/postgresql/9.6/main/postgresql.conf`
RUN echo "listen_addresses='*'" >> /etc/postgresql/9.6/main/postgresql.conf

RUN curl -sL https://storage.googleapis.com/golang/go1.16.5.linux-amd64.tar.gz | tar -C /usr/local -xz
RUN ln -s  /usr/local/go/bin/go /usr/local/bin/go
WORKDIR /go/src/go
COPY go.mod go.sum ./
RUN go mod download
COPY . ./

ENV PGPORT=5432
ENV PGUSER=circleci
ENV PGHOST=localhost

WORKDIR /go/src/go/exp/tools/dump-ledger-state

ARG GITCOMMIT
ENV GITCOMMIT=${GITCOMMIT}

ENTRYPOINT ["./docker-entrypoint.sh"]
