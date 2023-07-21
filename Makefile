COMMON_SELF_DIR := $(dir $(lastword $(MAKEFILE_LIST)))
MAKEFILE_DIR := ${COMMON_SELF_DIR}

ifeq ($(origin ROOT_DIR),undefined)
ROOT_DIR := $(MAKEFILE_DIR)
endif


.PHONY: grpc
# generate grpc code(no mi); grpc-{api version}
grpc:
	@echo "generate grpc..."
	@cd ${ROOT_DIR}api/v1 && protoc --proto_path=. \
	       --proto_path=../third_party \
           --go_out=paths=source_relative:.\
           --go-grpc_out=paths=source_relative:. \
           $(shell cd ${ROOT_DIR}api/v1 && find . -name "*.proto")
	@echo "generate grpc finsh."

.PHONY: build
build:
	@go build -o gid -ldflags "-w -s -X main.APPVersion=$(QM_APP_VERSION)" -tags=jsoniter cmd/main.go cmd/wire_gen.go

.PHONY: run-dev
run-dev:
	@go build -o gid-dev cmd/main.go cmd/wire_gen.go \
	&& ./gid-dev -conf ${ROOT_DIR}/configs/configs.yaml