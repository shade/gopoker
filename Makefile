PROTO_SRC_DIR=proto
PROTO_TARGET_DIR=$(GOPATH)/src/poker_backend/target

all:
	go run main.go -- -wsport=8081

build:
	protoc $(PROTO_SRC_DIR)/messages.proto --go_out=$(PROTO_TARGET_DIR)