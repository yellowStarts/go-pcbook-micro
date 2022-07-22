gen:
	protoc --proto_path=proto proto/*.proto --go_out=plugins=grpc:pb

clean:
	rm -rf pb/*

server:
	go run cmd/server/main.go -port 8080

client:
	go run cmd/client/main.go -address 0.0.0.0:8080

test:
# -cover 衡量测试的代码覆盖率 
# -race 检测代码中的 race 情况
	go test -cover -race ./...

cert: # 前提需要安装 openssl
	cd cert; ./gen.sh; cd ..

.PHONY: gen clean server client test cert