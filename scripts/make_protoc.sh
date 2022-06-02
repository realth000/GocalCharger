#!/bin/bash

which protoc 1>/dev/null 2>&1
[ $? != 0 ] && echo "protoc not found, exit" && exit 1

src_dir=`dirname $0`/../net/proto
cd ${src_dir}

protoc --go_out=. \
	--go_opt=paths=source_relative \
	--go-grpc_out=. \
	--go-grpc_opt=paths=source_relative \
	--proto_path=. \
	gocalcharger.proto
