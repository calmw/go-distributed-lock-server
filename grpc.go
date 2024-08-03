package main

import (
	"context"
	"distributed-lock/service"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	healthsvc "google.golang.org/grpc/health"
	healthz "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

const DefaultTimeout = 2

type LockServer struct {
	service.UnimplementedLockServiceServer
}

func registerServices(s *grpc.Server, h *healthsvc.Server) {
	service.RegisterLockServiceServer(s, &LockServer{})
	healthz.RegisterHealthServer(s, h)
	reflection.Register(s)
}

func startServer(s *grpc.Server, l net.Listener) error {
	return s.Serve(l)
}

func stopServer(s *grpc.Server, h *healthsvc.Server, d time.Duration) {
	updateServiceHealth(h, service.LockService_ServiceDesc.ServiceName, healthz.HealthCheckResponse_NOT_SERVING)
	time.Sleep(d)
	s.Stop()
	s.GracefulStop()
}

func updateServiceHealth(
	h *healthsvc.Server,
	service string,
	status healthz.HealthCheckResponse_ServingStatus,
) {
	h.SetServingStatus(service, status)
}

func (s *LockServer) Lock(ctx context.Context, req *service.LockRequest) (*service.LockReply, error) {
	res, msg := service.Lock(req.LockName, req.ClientId)
	return &service.LockReply{
		Result: res,
		Msg:    msg,
	}, nil
}

func (s *LockServer) UnLock(ctx context.Context, req *service.UnLockRequest) (*service.UnLockReply, error) {
	res, msg := service.UnLock(req.LockName, req.ClientId)
	return &service.UnLockReply{
		Result: res,
		Msg:    msg,
	}, nil
}

func (s *LockServer) ForceLock(ctx context.Context, req *service.ForceLockRequest) (*service.ForceLockReply, error) {
	res, msg := service.ForceLock(req.LockName, req.LockName)
	return &service.ForceLockReply{
		Result: res,
		Msg:    msg,
	}, nil
}

func (s *LockServer) ForceUnLock(ctx context.Context, req *service.ForceUnLockRequest) (*service.ForceUnLockReply, error) {
	res, msg := service.ForceUnLock(req.LockName)
	return &service.ForceUnLockReply{
		Result: res,
		Msg:    msg,
	}, nil
}

// 服务端，一元RPC方法调用的日志拦截器
func loggingUnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	logMessage(ctx, info.FullMethod, time.Since(start), err)

	return resp, err
}

// 服务端，一元紧急处理拦截器
func panicUnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic recovered: %v", r)
			err = status.Error(
				codes.Internal,
				"Unexpected error happened",
			)
		}
	}()
	resp, err = handler(ctx, req)

	return
}

// 服务端，一元超时终止请求拦截器
// 假设我们想对RPC方法执行时间施加一个上限，我们知道对于某些恶意请求用户，RPC调用方法可能需要比300毫秒更长的时间，这种情况下我们只想终止请求
// 任何超过300毫秒的RPC方法都将被终止
func timeoutUnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	var resp interface{}
	var err error

	timeoutStr := os.Getenv("TIME_OUT")
	timeout, err := strconv.Atoi(timeoutStr)
	if err != nil {
		timeout = DefaultTimeout
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	ch := make(chan error)
	go func() {
		resp, err = handler(ctxWithTimeout, req)
		ch <- err
	}()

	select {
	case <-ctxWithTimeout.Done():
		cancel()
		err = status.Error(
			codes.DeadlineExceeded,
			fmt.Sprintf("%s: DeadlineExceeded", info.FullMethod),
		)
		return resp, err
	case <-ch:
	}

	return resp, err
}

// 记录RPC方法调用的详细信息
func logMessage(
	ctx context.Context,
	method string,
	latency time.Duration,
	err error,
) {
	var client string
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Println("No metadata in context")
	} else {
		if len(md.Get("client")) != 0 {
			client = md.Get("client")[0]
		}
	}
	log.Printf(
		"Method: %s, Latency: %v, Error: %v, client: %s",
		method,
		latency,
		err,
		client,
	)
}

// 下面的一个结构体以及方法，是对服务端流的包装，将使用这些方法对原本流处理方法进行替换，来对服务端流的包装，实现每次流传输都可以进行自定义操作，而不是原本的等到全部传输完成才执行拦截器
type wrappedServerStream struct {
	RecvMsgTimeout time.Duration // 流超时时间
	grpc.ServerStream
}

func (s wrappedServerStream) SendMsg(m interface{}) error {
	log.Printf("Send msg called: %T", m)
	return s.ServerStream.SendMsg(m)
}

func (s wrappedServerStream) RecvMsg(m interface{}) error {
	ch := make(chan error)
	t := time.NewTimer(s.RecvMsgTimeout)
	go func() {
		log.Printf("Waiting to receive a msg: %T", m)
		ch <- s.ServerStream.RecvMsg(m)
	}()
	select {
	case <-t.C:
		return status.Error(
			codes.DeadlineExceeded,
			"Deadline exceeded",
		)
	case err := <-ch:
		return err
	}
}
