package main

import (
	"context"
	"distributed-lock/service"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"os"
)

const (
	DefaultServerAddr = "127.0.0.1:6000"
	ClientId          = "XX"
	LockName          = "order_lock"
)

func main() {
	serverAddr := os.Getenv("SERVER_ADDR")
	if len(serverAddr) == 0 {
		serverAddr = DefaultServerAddr
	}

	c := NewLockServiceClient(serverAddr)

	result, err := c.Lock(context.Background(), &service.LockRequest{
		ClientId: ClientId,
		LockName: LockName,
	})

	s := status.Convert(err) // status.Convert函数分别访问错误代码和错误消息
	if s.Code() != codes.OK {
		log.Fatalf("Request failed: %v-%v\n", s.Code(), s.Message())
	}
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)
}
