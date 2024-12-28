PROTO_DIR=proto
OUT_DIR=proto

PROTO_FILES=$(wildcard $(PROTO_DIR)/*.proto)

generate:
	protoc --proto_path=proto \
	       --go_out=proto --go_opt=paths=source_relative \
	       --go-grpc_out=proto --go-grpc_opt=paths=source_relative \
	       proto/*.proto

clean:
	rm -f $(OUT_DIR)/*.pb.go

.PHONY: generate clean
