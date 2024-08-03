package main

import (
	"distributed-lock/service"
	"google.golang.org/grpc"
	healthsvc "google.golang.org/grpc/health"
	healthz "google.golang.org/grpc/health/grpc_health_v1"
	"log"
	"net"
	"os"
)

func main() {
	service.InitLock()

	listenAddr := os.Getenv("LISTEN_ADDR")
	if len(listenAddr) == 0 {
		listenAddr = "0.0.0.0:6000"
	}
	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor( // 用于注册多个服务端一元拦截器，最内层的拦截器首先执行
			loggingUnaryInterceptor,
			timeoutUnaryInterceptor,
			panicUnaryInterceptor, //
			// ... 其他拦截器
		),
	)

	h := healthsvc.NewServer()
	registerServices(s, h)
	updateServiceHealth(h, service.LockService_ServiceDesc.ServiceName, healthz.HealthCheckResponse_SERVING)
	log.Fatal(startServer(s, lis))
}
