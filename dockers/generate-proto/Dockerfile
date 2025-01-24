# Use Ubuntu as base image as this image is just to generate the protobug bindings,
# not to run a service.
FROM ubuntu
#Install protocol buffers
RUN apt-get update
RUN apt install -y build-essential
RUN apt install -y protobuf-compiler
RUN apt install -y protoc-gen-go
RUN apt install -y curl
RUN curl https://sh.rustup.rs -sSf | bash -s -- -y
RUN echo 'source $HOME/.cargo/env' >> $HOME/.bashrc
RUN bash -c "source /root/.cargo/env; cargo install protobuf-codegen"

# Add the generating script
WORKDIR /home
COPY protoc.sh /home/protoc.sh
ENTRYPOINT ["/bin/bash","-c","/home/protoc.sh"]
#ENTRYPOINT ["/bin/bash"]
