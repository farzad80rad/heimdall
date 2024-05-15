#!/bin/sh
 protoc -I . \
  --grpc-gateway_out ../Proto/gen/Proto \
  --go_out ../Proto/gen/Proto --go_opt paths=source_relative \
  --go-grpc_out ../Proto/gen/Proto --go-grpc_opt paths=source_relative \
  --grpc-gateway_opt paths=source_relative \
  --grpc-gateway_opt generate_unbound_methods=true \
   ./*.proto


