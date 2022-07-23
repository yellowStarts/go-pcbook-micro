gen:
	protoc --proto_path=proto proto/*.proto --go_out=plugins=grpc:pb

clean:
	rm -rf pb/*

server1:
	go run cmd/server/main.go -port 50051

server2:
	go run cmd/server/main.go -port 50052

server1-tls:
	go run cmd/server/main.go -port 50051 -tls

server2-tls:
	go run cmd/server/main.go -port 50052 -tls

server:
	go run cmd/server/main.go -port 8080

client:
	go run cmd/client/main.go -address 127.0.0.1:8080

client-tls:
	go run cmd/client/main.go -address 127.0.0.1:8080 -tls

test:
# -cover 衡量测试的代码覆盖率 
# -race 检测代码中的 race 情况
	go test -cover -race ./...

cert: # 前提需要安装 openssl
	cd cert; ./gen.sh; cd ..

.PHONY: gen clean server1 server2 server client test cert 