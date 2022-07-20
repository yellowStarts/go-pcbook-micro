# gRPC 学习案例

## 安装
1. protoc 安装

- 进入 [protobuf release](https://github.com/protocolbuffers/protobuf/releases) 页面，选择适合自己操作系统的压缩包文件
- 解压 `protoc-x.x.x-osx-x86_64.zip` 并进入 `protoc-x.x.x-osx-x86_64`
```
$ cd protoc-x.x.x-osx-x86_64/bin
```
- 将启动的 `protoc` 二进制文件移动到被添加到环境变量的任意 path下，如 `$GOPATH/bin`，这里不建议直接将其和系统的以下path放在一起。
```
$ mv protoc $GOPATH/bin
```
>tip: `$GOPATH`为你本机的实际文件夹地址
- 验证安装结果
```
$ protoc --version
libprotoc x.x.x
```

2. protoc-gen-go/protoc-gen-go-grpc 安装
```
$ go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
$ go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

3. gRPC 资料
- 链接：[gRPC基于Go快速开始](https://grpc.io/docs/languages/go/quickstart/)
- 安装：`go get -u google.golang.org/grpc`

## protoc 的使用
`protoc` 生成 GO 代码：
```
$ protoc --proto_path=proto proto/*.proto --go_out=plugins=grpc:pb
```