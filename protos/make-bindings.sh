#!/usr/bin/env bash
# Use the protoc image to run protoc.sh and generate the bindings.
docker run -e PROTO=request.proto --mount type=bind,source="$PWD",target=/home/proto/ -it saichler/protoc:latest
docker run -e PROTO=priority.proto --mount type=bind,source="$PWD",target=/home/proto/ -it saichler/protoc:latest
docker run -e PROTO=message.proto --mount type=bind,source="$PWD",target=/home/proto/ -it saichler/protoc:latest
docker run -e PROTO=tests.proto --mount type=bind,source="$PWD",target=/home/proto/ -it saichler/protoc:latest

# Now move the generated bindings to the models directory and clean up
mkdir -p ../go/types
mv ./types/*.pb.go ../go/types/.
rm -rf ./types

mkdir -p ../go/tests
mv ./tests/*.pb.go ../go/tests/.
rm -rf ./tests