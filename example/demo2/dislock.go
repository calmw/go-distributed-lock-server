package main

import (
	"context"
	"distributed-lock/service"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"os"
)

const (
	LockServerAddr = "127.0.0.1:6000"
	LockName       = "order_lock"
)

var Lk service.LockServiceClient

func InitLock() {
	serverAddr := os.Getenv("LOCK_SERVER_ADDR")
	if len(serverAddr) == 0 {
		serverAddr = LockServerAddr
	}

	Lk = NewLockServiceClient(serverAddr)
}

func Lock(clientId string) (bool, error) {
	result, err := Lk.Lock(context.Background(), &service.LockRequest{
		ClientId: clientId,
		LockName: LockName,
	})

	s := status.Convert(err) // status.Convert函数分别访问错误代码和错误消息
	if s.Code() != codes.OK || err != nil {
		log.Fatalf("Request failed: %v-%v\n", s.Code(), s.Message())
		return false, err
	}
	return result.Result, errors.New(result.Msg)
}

func UnLock(clientId string) (bool, error) {
	result, err := Lk.UnLock(context.Background(), &service.UnLockRequest{
		ClientId: clientId,
		LockName: LockName,
	})

	s := status.Convert(err) // status.Convert函数分别访问错误代码和错误消息
	if s.Code() != codes.OK || err != nil {
		log.Fatalf("Request failed: %v-%v\n", s.Code(), s.Message())
		return false, err
	}
	return result.Result, errors.New(result.Msg)
}

func ForceUnLock() (bool, error) {
	result, err := Lk.ForceUnLock(context.Background(), &service.ForceUnLockRequest{
		LockName: LockName,
	})

	s := status.Convert(err) // status.Convert函数分别访问错误代码和错误消息
	if s.Code() != codes.OK || err != nil {
		log.Fatalf("Request failed: %v-%v\n", s.Code(), s.Message())
		return false, err
	}
	return result.Result, errors.New(result.Msg)
}
