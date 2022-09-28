.PHONY: protos

protos:
	protoc --proto_path=protos --go_out=protos/currency --go_opt=paths=source_relative protos/currency.proto

# https://developers.google.com/protocol-buffers/docs/reference/go-generated