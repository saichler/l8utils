# Use Ubuntu as base image as this image is just to generate the protobug bindings,
# not to run a service.
FROM ubuntu

#Install protocol buffers
RUN apt-get update
RUN apt install -y build-essential
RUN apt install -y protobuf-compiler
RUN apt install -y curl
# Install protoc-gen-go pinned to v1.32.0 (do NOT use apt's protoc-gen-go, which tracks
# whatever version is in the Ubuntu repos and drifts to newer releases over time)
RUN curl -LO https://go.dev/dl/go1.22.2.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.22.2.linux-amd64.tar.gz && \
    rm go1.22.2.linux-amd64.tar.gz
ENV PATH="/usr/local/go/bin:/root/go/bin:${PATH}"
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.32.0
RUN cp /root/go/bin/protoc-gen-go /usr/local/bin/protoc-gen-go && chmod 755 /usr/local/bin/protoc-gen-go
RUN curl https://sh.rustup.rs -sSf | bash -s -- -y
RUN echo 'source $HOME/.cargo/env' >> $HOME/.bashrc
RUN bash -c "source /root/.cargo/env; cargo install protobuf-codegen"
# Make protoc-gen-rs available to non-root users (needed when docker run --user ... is used)
RUN cp /root/.cargo/bin/protoc-gen-rs /usr/local/bin/protoc-gen-rs && chmod 755 /usr/local/bin/protoc-gen-rs
# Add the generating script
WORKDIR /home
COPY protoc.sh /home/protoc.sh
ENTRYPOINT ["/bin/bash","-c","/home/protoc.sh"]
#ENTRYPOINT ["/bin/bash"]
