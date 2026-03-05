#!/bin/sh
DBNAME=$1
DBUSER=$2
DBPASSWORD=$3
DBPORT=${4:-5432}
if [ -z "$DBNAME" ]; then echo "Error: DBNAME (arg 1) is required"; exit 1; fi
if [ -z "$DBUSER" ]; then echo "Error: DBUSER (arg 2) is required"; exit 1; fi
if [ -z "$DBPASSWORD" ]; then echo "Error: DBPASSWORD (arg 3) is required"; exit 1; fi
DBDIR="/data/postgres/$DBNAME"
mkdir -p /run/postgresql
mkdir -p /data/postgres
if [ ! -f "$DBDIR/PG_VERSION" ]; then
  echo "Initializing database on port $DBPORT..."
  mkdir -p "$DBDIR"
  chmod 0700 "$DBDIR"
  initdb -D $DBDIR
  echo "listen_addresses = '*'" >> $DBDIR/postgresql.conf
  echo "port = $DBPORT" >> $DBDIR/postgresql.conf
  echo "shared_preload_libraries = 'timescaledb'" >> $DBDIR/postgresql.conf
  cat > $DBDIR/pg_hba.conf << EOF
local   all             postgres                                peer
host    all             all             0.0.0.0/0               md5
host    all             all             ::1/128                 md5
EOF
  pg_ctl -D $DBDIR -l /var/log/postgresql/postgresql.log start
  sleep 2
  psql -p $DBPORT -c "CREATE USER $DBUSER WITH PASSWORD '$DBPASSWORD';"
  psql -p $DBPORT -c "CREATE DATABASE $DBNAME OWNER $DBUSER;"
  psql -p $DBPORT -c "GRANT ALL PRIVILEGES ON DATABASE $DBNAME TO $DBUSER;"
  psql -p $DBPORT -d $DBNAME -c "CREATE EXTENSION IF NOT EXISTS timescaledb;"
  echo "Database initialized."
else
  echo "Database already exists, starting on port $DBPORT..."
  sed -i "s/^port = .*/port = $DBPORT/" $DBDIR/postgresql.conf
  pg_ctl -D $DBDIR -l /var/log/postgresql/postgresql.log start
fi
