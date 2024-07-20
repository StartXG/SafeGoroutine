#generate proto files
protoc --go_out=./proto ./proto/*.proto && protoc --go-grpc_out=./proto ./proto/*.proto

