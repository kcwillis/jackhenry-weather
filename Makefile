
PROTO_PATH_V1=proto/v1
PROTO_PATH_GOOGLEAPIS=${GOPATH}/src/github.com/googleapis/googleapis

dep:
	go get -u \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
    google.golang.org/protobuf/cmd/protoc-gen-go \
    google.golang.org/grpc/cmd/protoc-gen-go-grpc

mod:
	go mod tidy
	go mod vendor

protoc:
	protoc ${PROTO_PATH_V1}/*.proto \
	--proto_path=${PROTO_PATH_V1} \
	--proto_path=${PROTO_PATH_GOOGLEAPIS} \
	--go_out=${PROTO_PATH_V1} \
	--go_opt=paths=source_relative \
	--go-grpc_out=${PROTO_PATH_V1} \
	--go-grpc_opt=paths=source_relative \
	--grpc-gateway_out ${PROTO_PATH_V1} --grpc-gateway_opt paths=source_relative

build:
	@echo TODO: build

helm:
	@echo TODO: helm

container:
	@echo TODO: build docker container
