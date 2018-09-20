.PHONY: all

all:
	@/home/udo/Downloads/protoc-3.6.1-linux-x86_64/bin/protoc --proto_path=definitions --go_out=plugins=grpc:. global.proto 2fa.proto admin.proto db.proto jwt.proto mail.proto pay.proto
	@echo "Protofiles rebuild done"