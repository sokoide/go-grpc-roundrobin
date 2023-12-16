

.PHONY: all readme docker server client proto install clean

all: docker server client

readme:
	echo Take a look at https://grpc.io/docs/languages/go/quickstart/

docker: server client
	docker build -t grpc .

server:
	go build ./cmd/server

client:
	go build ./cmd/client

proto: proto/hello.proto
	protoc --proto_path=proto --go_out=./pkg/grpc --go-grpc_out=. --go_opt=paths=source_relative hello.proto

install:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

clean:
	go clean
