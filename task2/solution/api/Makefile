OUT = gen/.
PROTOS = $(shell find . -type f -name '*.proto')

generate:
	protoc \
		--go_out=$(OUT) \
		--go_opt=paths=source_relative \
		--go-grpc_out=$(OUT) \
		--go-grpc_opt=paths=source_relative \
		$(PROTOS)
