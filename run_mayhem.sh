#!/bin/bash

set -ex

HOST=localhost:8081
PROTOCOL=http

mapi run \
     training/mapi-grpc-example 30 example-api.swagger.json \
     --url "$PROTOCOL://$HOST/" \
     --resource-hint "POST /mapi_grpc_example.api.v1.UserService/GetUsers BODY source:valid-source" \
     --resource-hint "POST /mapi_grpc_example.api.v1.UserService/CheckReservedName BODY source:valid-source" \
     --interactive \
     "${@}" 
