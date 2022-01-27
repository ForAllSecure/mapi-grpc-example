# Generate code from .proto files

# Dependencies (tools must be available on the path):
# - protoc (Protocol Buffer v3)
# - protoc-gen-go (protoc compiler extension for Go)

######################################################
# Run all code generation tasks
######################################################
all: proto_go_stub proto_proxy_stub swagger

######################################################
# Clean all generated code
######################################################
clean:
	rm -f *.pb.go *.pb.gw.go *.swagger.json

######################################################
# Generate protobuf stubs
# See https://github.com/golang/protobuf)
######################################################
proto_go_stub: example-api.pb.go
example-api.pb.go: *.proto
	protoc -I . *.proto --go-grpc_out=. --go_out=.

######################################################
# Generate protobuf proxy stubs
# See https://github.com/grpc-ecosystem/grpc-gateway
######################################################
proto_proxy_stub:
	protoc -I . --grpc-gateway_out . \
	--grpc-gateway_opt logtostderr=true \
	--grpc-gateway_opt paths=source_relative \
	--grpc-gateway_opt generate_unbound_methods=true \
	*.proto

######################################################
# Generate swagger (OpenAPI) spec from .proto
######################################################
swagger: example-api.swagger.json
example-api.swagger.json:
	protoc -I . --openapiv2_out . \
	--openapiv2_opt logtostderr=true \
	--openapiv2_opt generate_unbound_methods=true \
	*.proto