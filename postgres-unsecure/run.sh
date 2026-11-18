#sudo pkill -f postgres
#sudo rm -rf /data/postgres/problerdb
#sudo rm -rf /data/postgres/probleralarms
docker run -d -p 5432:5432 -v /data/:/data/ registry.layer-8.dev:44443/library/unsecure-postgres:latest admin admin admin 5432
docker run -d -p 5433:5433 -v /data/:/data/ registry.layer-8.dev:44443/library/unsecure-postgres:latest admin admin admin 5433
