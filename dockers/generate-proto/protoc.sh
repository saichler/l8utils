#!/usr/bin/env bash
echo $(uname -a)
echo "Generating Protobufs for $PROTO"
cd proto
protoc --go_out=. $PROTO
protoc --rs_out=. $PROTO
echo "Done!"