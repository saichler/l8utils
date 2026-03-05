ARG NAME='alpine:latest'
FROM $NAME AS final

ARG HOST_UID=1000
ARG HOST_GID=1000

RUN apk add postgresql postgresql-contrib

RUN sed -i "s/^postgres:x:[0-9]*:[0-9]*/postgres:x:${HOST_UID}:${HOST_GID}/" /etc/passwd && \
    sed -i "s/^postgres:x:[0-9]*/postgres:x:${HOST_GID}/" /etc/group

# Build and install TimescaleDB from source
RUN apk add --no-cache --virtual .build-deps \
        git cmake build-base openssl-dev postgresql-dev python3 bash
RUN git clone --branch 2.23.0 --depth 1 \
        https://github.com/timescale/timescaledb.git /tmp/timescaledb && \
    cd /tmp/timescaledb && \
    ./bootstrap -DREGRESS_CHECKS=OFF -DWARNINGS_AS_ERRORS=OFF -DAPACHE_ONLY=1 && \
    cd build && make && make install && \
    cd / && rm -rf /tmp/timescaledb
RUN apk del .build-deps
RUN mkdir -p /run/postgresql
RUN chown -R postgres:postgres /run/postgresql/
RUN mkdir -p /data/postgres
RUN chown -R postgres:postgres /data/postgres

EXPOSE 5432
EXPOSE 5433

RUN mkdir -p /var/log/postgresql && \
    chown postgres:postgres /var/log/postgresql

COPY start-postgres.sh /start-postgres.sh
RUN chmod +x /start-postgres.sh

USER postgres

ENTRYPOINT ["/bin/sh", "-c", "/start-postgres.sh $0 $1 $2 $3 && sleep 3600"]
