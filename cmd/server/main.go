package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"go-pcbook-micro/pb"
	"go-pcbook-micro/service"
	"io/ioutil"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

func seedUsers(userStore service.UserStore) error {
	err := createUser(userStore, "admin1", "secret", "admin")
	if err != nil {
		return err
	}

	return createUser(userStore, "user1", "secret", "user")
}

func createUser(userStore service.UserStore, username, password, role string) error {
	user, err := service.NewUser(username, password, role)
	if err != nil {
		return err
	}
	return userStore.Save(user)
}

const (
	secretKey     = "secret"
	tokenDuration = 15 * time.Minute
)

func accessibleRoles() map[string][]string {
	const laptopServicePath = "/proto.LaptopService/"
	return map[string][]string{
		laptopServicePath + "CreateLaptop": {"admin"},
		laptopServicePath + "UploadImage":  {"admin"},
		laptopServicePath + "RateLaptop":   {"admin", "user"},
	}
}

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	// 加载签署客户端证书的CA的证书
	pemClientCA, err := ioutil.ReadFile("cert/ca-cert.pem")
	if err != nil {
		return nil, err
	}
	// 创建 x509证书池
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemClientCA) {
		return nil, fmt.Errorf("failed to add client CA's certificate")
	}
	// 加载服务器证书和私钥
	serverCert, err := tls.LoadX509KeyPair("cert/server-cert.pem", "cert/server-key.pem")
	if err != nil {
		return nil, err
	}
	// 创建 tls 配置对象
	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	}

	return credentials.NewTLS(config), nil
}

func main() {
	port := flag.Int("port", 0, "the server port")
	enableTLS := flag.Bool("tls", false, "enable SSL/TLS")
	flag.Parse()
	log.Printf("start server on port %d, TLS = %t", *port, *enableTLS)

	// authService
	userStore := service.NewInMemoryUserStore()
	// 创建种子用户
	err := seedUsers(userStore)
	if err != nil {
		log.Fatal("cannot seed users")
	}
	jwtManager := service.NewJWTManager(secretKey, tokenDuration)
	authService := service.NewAuthService(userStore, jwtManager)
	// laptopServer
	laptopStore := service.NewInMemoryLaptopStore()
	imageStore := service.NewDiskImageStore("images")
	ratingStore := service.NewInMemoryRatingStore()
	laptopServer := service.NewLaptopServer(laptopStore, imageStore, ratingStore)

	// 拦截器
	interceptor := service.NewAuthInterceptor(jwtManager, accessibleRoles())
	serverOptions := []grpc.ServerOption{
		grpc.UnaryInterceptor(interceptor.Unary()),   // 一元rpc拦截器
		grpc.StreamInterceptor(interceptor.Stream()), // 流式rpc拦截器
	}
	if *enableTLS {
		// 获取tls 凭据对象
		tlsCredentials, err := loadTLSCredentials()
		if err != nil {
			log.Fatal("cannot load TLS credentials: ", err)
		}

		serverOptions = append(serverOptions, grpc.Creds(tlsCredentials))
	}
	// grpcServer
	grpcServer := grpc.NewServer(serverOptions...)
	// 注册服务
	pb.RegisterAuthServiceServer(grpcServer, authService)
	pb.RegisterLaptopServiceServer(grpcServer, laptopServer)
	// 反射
	reflection.Register(grpcServer)

	address := fmt.Sprintf("0.0.0.0:%d", *port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("canont start server: ", err)
	}
}
