# Use Ubuntu as base image as this image is just to generate the protobug bindings,
# not to run a service.
FROM ubuntu
#Install protocol buffers
RUN apt-get update
RUN apt install -y protobuf-compiler
RUN apt install -y protoc-gen-go
# Add the generating script
WORKDIR /home
COPY ./protoc.sh /home/protoc.sh

ENTRYPOINT ["/bin/bash","-c","/home/protoc.sh"]
