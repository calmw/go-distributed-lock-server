package main

import (
	"context"
	"distributed-lock/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
)

func NewLockServiceClient(serverAddr string) service.LockServiceClient {
	log.Println("connecting grpc server")
	conn, err := setupGrpcConnection(serverAddr)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("grpc server connected")
	//return svc.NewScannerServiceClient(conn)
	return service.NewLockServiceClient(conn)
}

func setupGrpcConnection(addr string) (*grpc.ClientConn, error) {
	return grpc.DialContext(
		context.Background(),
		addr,
		grpc.WithBlock(),    // 确保在函数返回之前建立连接。这意味着如果在服务器启动并运行之前运行客户端，它将无限期等待。
		grpc.WithInsecure(), // 用于注册多个客户端一元拦截器，最内层的拦截器首先执行
		grpc.WithChainUnaryInterceptor( // 用于注册多个客户端一元拦截器，最内层的拦截器首先执行
			metadataUnaryInterceptor,
			// ... 其他拦截器
		),
	)
}

// 一元客户端拦截器，对传出的任何一元RPC请求添加一个唯一标识符，添加的数据在context中
func metadataUnaryInterceptor(
	ctx context.Context,
	method string,
	req, reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	ctxWithMetadata := metadata.AppendToOutgoingContext(
		ctx,
		"client",
		"node",
	)
	return invoker(
		ctxWithMetadata,
		method,
		req,
		reply,
		cc,
		opts...,
	)
}
