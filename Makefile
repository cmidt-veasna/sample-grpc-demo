# docker configure
TAG ?= 1.0.0

# protobuf & grpc configure
# protoc binary
PROTOC ?= protoc
# current directory
CUR_DIR ?= $(shell pwd)
# include proto directory
PROTO_INCLUDE ?= $(CUR_DIR)
# generate code parent directory
GEN_DIR ?= $(PROTO_INCLUDE)/server/go/pkg/example
# proto flag
FLAGS += -I$(PROTO_INCLUDE)/proto
FLAGS += -I$(PROTO_INCLUDE)/proto/googleapis
# final flag
DESCRIPTOR_FLAGS = $(FLAGS) --include_imports --include_source_info
GRPC_FLAGS = $(FLAGS) --go_out=plugins=grpc:$(GEN_DIR)
SWAGGER_FLAGS = $(FLAGS) --swagger_out=logtostderr=true:$(CUR_DIR)/clients/web/rest/

# client code generate
MAKE ?= make
CLIENT ?= go
CLIENT_GO_FLAGS = $(FLAGS) --go_out=plugins=grpc:$(CUR_DIR)/clients/go/pkg/example
CLIENT_IOS_FLAGS = $(FLAGS) --swift_out='$(CUR_DIR)/clients/ios/example/Example GrpcTests' --swiftgrpc_out='$(CUR_DIR)/clients/ios/example/Example GrpcTests'

clean:
	@echo "clean generate code ...."
	@rm -f $(CUR_DIR)/dockers/*.pb && rm -f $(GEN_DIR)/*.pb.go

generate: clean
	@echo "Generate go server side stub ...."
	@$(PROTOC) $(GRPC_FLAGS) $(CUR_DIR)/proto/*.proto

	@echo "Generate grpc transcoder ...."
	@$(PROTOC) $(DESCRIPTOR_FLAGS) --descriptor_set_out=$(CUR_DIR)/proto/descriptor.pb $(CUR_DIR)/proto/*.proto
	@mv $(CUR_DIR)/proto/*.pb dockers/

build: generate
	@echo "Build Server ..."
	@cd server/go && GOOS=linux GOARCH=amd64 go build -o server && mv server ../../dockers

	@echo "Build Docker image ..."
	@docker build -t envoy-grpc:$(TAG) dockers

	@echo "Clean up ..."
	@rm dockers/server && rm dockers/*.pb

generate-client:
ifeq ($(CLIENT),go)
generate-client: generate-client-go
else ifeq ($(CLIENT),java)
generate-client: generate-client-java
else ifdef ($(CLIENT),ios)
generate-client: generate-client-ios
else ifdef ($(CLIENT),web)
generate-client: generate-client-web
endif

clean-client: generate-client-go-clean

generate-client-go: generate-client-go-clean
	@echo "Generate client code $(CLIENT) ..."
	@$(PROTOC) $(CLIENT_GO_FLAGS) $(CUR_DIR)/proto/*.proto

	@echo "Generate go client ..."
	@cd clients/go && go build -o sample

generate-client-go-clean:
	@echo "Clean go client code ..."
	@rm -f $(CUR_DIR)/clients/go/pkg/example/*.pb.go
	@rm -f $(CUR_DIR)/clients/go/sample

generate-client-java:
	@cd clients/java && $(MAKE) clean && $(MAKE) build

generate-client-ios:
	@echo "Clean generated code ..."

	@echo "Generate swift code ..."
	@$(PROTOC) $(CLIENT_IOS_FLAGS) $(CUR_DIR)/proto/*.proto

generate-client-web:
	@echo "Clean ..."
	@rm $(CUR_DIR)/clients/web/rest/*.swagger.json

	@echo "Generate Swagger file ...."
	@$(PROTOC) $(SWAGGER_FLAGS) $(CUR_DIR)/proto/*.proto

generate-client-web-compile:
	@echo "Compile swagger js client ..."
	@cd $(CUR_DIR)/clients/web && $(MAKE) clean && $(MAKE) build
