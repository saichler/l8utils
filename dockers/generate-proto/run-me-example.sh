#!/usr/bin/env bash
# Use the protoc image to run protoc.sh and generate the bindings.
docker run --user "$(id -u):$(id -g)" -e PROTO=my-model.proto --mount type=bind,source="$PWD",target=/home/proto -it saichler/protoc:latest

# Now move the generated bindings to the models directory and clean up
mkdir -p ../go/models
mv ./models/my-model.pb.go ../go/models/.
rm -rf ./models
