package main

import (
	"SafeGoroutine/proto"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"sync"
)

const PORT = ":8000"

var (
	balance int32 = 1000
	wg      sync.WaitGroup
	rwLock  sync.RWMutex
)

type server struct {
	proto.UnimplementedBankServiceServer
}

func (s *server) GetBalance(ctx context.Context, req *proto.Empty) (b *proto.Balance, e error) {
	// defer wg.Done()
	rwLock.RLock()
	b = &proto.Balance{BalanceNumber: balance}
	e = nil
	defer rwLock.RUnlock()
	return
}

func (s *server) ModifyNumber(ctx context.Context, req *proto.Action) (b *proto.Balance, e error) {
	rwLock.Lock()
	defer rwLock.Unlock()

	if balance+req.ActionNumber < 0 {
		log.Println("执行交易：", req.ActionNumber, "，现有余额：", balance, "，余额不足，禁止执行")
		return nil, fmt.Errorf("余额不足，无法执行操作")
	}

	log.Println("执行交易：", req.ActionNumber, "，现有余额：", balance, "，允许执行")
	balance += req.ActionNumber
	b = &proto.Balance{BalanceNumber: balance}
	return b, nil
}

func main() {
	lis, err := net.Listen("tcp", PORT)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	proto.RegisterBankServiceServer(s, &server{})
	reflection.Register(s)
	log.Printf("Server is running on port: %v", PORT)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
