protoc --proto_path=definitions --go_out=plugins=grpc:. global.proto 2fa.proto db.proto jwt.proto mail.proto notification.proto pay.proto
