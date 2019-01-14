.PHONY: all

all:
	@protoc --proto_path=definitions --go_out=plugins=grpc:. global.proto 2fa.proto admin.proto db.proto jwt.proto mail.proto pay.proto sse.proto
	@echo "Protofiles rebuild done"