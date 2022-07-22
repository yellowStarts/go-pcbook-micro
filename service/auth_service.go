package service

import (
	"context"
	"go-pcbook-micro/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AuthService 授权服务
type AuthService struct {
	userStore  UserStore
	jwtManager *JWTManager
}

// NewAuthService 创建授权服务实例
func NewAuthService(userStore UserStore, jwtManager *JWTManager) *AuthService {
	return &AuthService{userStore, jwtManager}
}

func (server *AuthService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	user, err := server.userStore.Find(req.Username)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot find user: %v", err)
	}
	if user == nil || !user.IsCorrectPassword(req.GetPassword()) {
		return nil, status.Errorf(codes.NotFound, "incorrect username/password")
	}

	token, err := server.jwtManager.Generate(user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot generate access token")
	}

	res := &pb.LoginResponse{
		AccessToken: token,
	}
	return res, nil
}
