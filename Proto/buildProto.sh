#!/bin/sh
 protoc -I . \
  --grpc-gateway_out ./gen/Proto \
  --go_out ./gen/Proto --go_opt paths=source_relative \
  --go-grpc_out ./gen/Proto --go-grpc_opt paths=source_relative \
  --grpc-gateway_opt paths=source_relative \
  --grpc-gateway_opt generate_unbound_methods=true \
   ./*.proto


