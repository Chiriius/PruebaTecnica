.PHONY: generate
generate:
	protoc --go_out=. --go-grpc_out=. api/transports/grpc/events.proto