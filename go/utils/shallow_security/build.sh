go build -buildmode=plugin -o ./loader.so Loader.go Provider.go
sudo mv ./loader.so /var/loader.so
