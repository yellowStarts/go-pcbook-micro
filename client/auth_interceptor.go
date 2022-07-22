package client

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// AuthInterceptor 授权拦截器
type AuthInterceptor struct {
	authClient  *AuthClient
	authMethods map[string]bool
	accessToken string
}

// NewAuthInterceptor 创建实例
func NewAuthInterceptor(authClient *AuthClient, authMethods map[string]bool, refreshDuration time.Duration) (*AuthInterceptor, error) {
	interceptor := &AuthInterceptor{
		authClient:  authClient,
		authMethods: authMethods,
	}

	err := interceptor.scheduleRefreshToken(refreshDuration)
	if err != nil {
		return nil, err
	}

	return interceptor, nil
}

func (ai *AuthInterceptor) scheduleRefreshToken(refreshDuration time.Duration) error {
	err := ai.refreshToken()
	if err != nil {
		return err
	}

	go func() {
		wait := refreshDuration
		for {
			time.Sleep(wait)
			err := ai.refreshToken()
			if err != nil {
				wait = time.Second
			} else {
				wait = refreshDuration
			}
		}
	}()
	return nil
}

func (ai *AuthInterceptor) refreshToken() error {
	accessToken, err := ai.authClient.Login()
	if err != nil {
		return err
	}
	ai.accessToken = accessToken
	log.Printf("token refresh: %v", accessToken)

	return nil
}

// Unary 创建一元rpc拦截器
func (ai *AuthInterceptor) Unary() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		log.Println("--> unary interceptor: ", method)

		if ai.authMethods[method] {
			return invoker(ai.attachToken(ctx), method, req, reply, cc, opts...)
		}
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// token附加到上下文
func (ai *AuthInterceptor) attachToken(ctx context.Context) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "authorization", ai.accessToken)
}

// Stream 创建 流式rpc拦截器
func (ai *AuthInterceptor) Stream() grpc.StreamClientInterceptor {
	return func(
		ctx context.Context,
		desc *grpc.StreamDesc,
		cc *grpc.ClientConn,
		method string,
		streamer grpc.Streamer,
		opts ...grpc.CallOption,

	) (grpc.ClientStream, error) {
		log.Println("--> stream interceptor: ", method)

		if ai.authMethods[method] {
			return streamer(ai.attachToken(ctx), desc, cc, method, opts...)
		}
		return streamer(ctx, desc, cc, method, opts...)
	}
}
