#!/bin/bash

set -ex

HOST=localhost:8081
PROTOCOL=http

mapi run \
     training/mapi-grpc-example 30 api/v1/example-api.swagger.json \
     --url "$PROTOCOL://$HOST/" \
     --interactive \
     "${@}" 
