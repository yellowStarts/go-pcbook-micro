package client

import (
	"context"
	"go-pcbook-micro/pb"
	"time"

	"google.golang.org/grpc"
)

//AuthClient 授权客户端
type AuthClient struct {
	service  pb.AuthServiceClient
	username string
	password string
}

//NewAuthClient  创建授权客户端
func NewAuthClient(cc *grpc.ClientConn, username string, password string) *AuthClient {
	service := pb.NewAuthServiceClient(cc)
	return &AuthClient{service, username, password}
}

// Login 登录
func (client *AuthClient) Login() (string, error) {
	// 设置超时
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.LoginRequest{
		Username: client.username,
		Password: client.password,
	}

	res, err := client.service.Login(ctx, req)
	if err != nil {
		return "", err
	}

	return res.GetAccessToken(), nil
}
