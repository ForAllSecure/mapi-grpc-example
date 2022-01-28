# Builds the gRPC server fronted by a REST API gateway
#
# See also
#  - Mayhem for API: https://mayhem4api.forallsecure.com/
#  - grpc-gateyway: https://github.com/grpc-ecosystem/grpc-gateway
#
FROM golang:1.17.6

# Install protoc, protocol buffer compiler
RUN apt update && \
    apt install -y protobuf-compiler \
    # supervisor allows launching the gRPC server AND rest proxy
                   supervisor && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /opt/grpc-example

# Install required go modules
COPY go.mod go.sum tools.go ./
RUN go install \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
    google.golang.org/protobuf/cmd/protoc-gen-go \
    google.golang.org/grpc/cmd/protoc-gen-go-grpc

# Copy the rest of the code
COPY . .

# Generate code from .proto and
# proxy stubs from .proto files
RUN make api-stubs

# Build the gRPC server
WORKDIR /opt/grpc-example/server
RUN go build

# Build the rest proxy
WORKDIR /opt/grpc-example/mapi
RUN go build -o rest-proxy

# Generate the swagger specification
WORKDIR /opt/grpc-example/api/v1
RUN protoc -I . --openapiv2_out . \
    --openapiv2_opt logtostderr=true \
    --openapiv2_opt generate_unbound_methods=true \
    example-api.proto

# 50051 - gRPC server
# 8081  - rest proxy
EXPOSE 50051 8081

# Configure supervisor so that the server and proxy are run
# with supervisor
ADD mapi/supervisord.conf /etc/supervisor/conf.d/supervisord.conf

# Working directory is the server folder
WORKDIR /opt/grpc-example/server

# Launch with supervisor
ENTRYPOINT ["/usr/bin/supervisord"]