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
	lock    sync.Mutex
)

type server struct {
	proto.UnimplementedBankServiceServer
}

func (s *server) ModifyNumber(ctx context.Context, req *proto.Action) (b *proto.Balance, e error) {
	lock.Lock()
	if balance+req.ActionNumber < 0 {
		log.Println("执行交易：", req.ActionNumber, "，现有余额：", balance, "，余额不足，禁止执行")
		e = fmt.Errorf("余额不足，无法执行操作")
	} else {
		log.Println("执行交易：", req.ActionNumber, "，现有余额：", balance, "，允许执行")
		balance += req.ActionNumber
		e = nil
	}
	b = &proto.Balance{BalanceNumber: balance}
	defer lock.Unlock()
	return
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
