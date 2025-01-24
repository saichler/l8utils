#!/usr/bin/env bash
echo $(uname -a)
echo "Generating Protobufs for $PROTO"
cd proto
source /root/.cargo/env
protoc --go_out=. $PROTO
protoc --rs_out=. $PROTO
echo "Done!"