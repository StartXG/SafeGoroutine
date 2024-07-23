package main

import (
	"SafeGoroutine/proto"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"sync/atomic"
)

const PORT = ":8000"

var (
	balance int32 = 1000
)

type server struct {
	proto.UnimplementedBankServiceServer
}

func (s *server) ModifyNumber(ctx context.Context, req *proto.Action) (b *proto.Balance, e error) {
	// 先使用atomic.LoadInt32来获取当前的balance
	currentBalance := atomic.LoadInt32(&balance)

	// 计算新的余额
	newBalance := currentBalance + req.ActionNumber

	// 检查余额是否足够
	if newBalance < 0 {
		log.Println("执行交易：", req.ActionNumber, "，现有余额：", currentBalance, "，余额不足，禁止执行")
		e = fmt.Errorf("余额不足，无法执行操作")
		b = &proto.Balance{BalanceNumber: currentBalance}
		return
	}

	// 使用atomic.AddInt32来更新balance
	atomic.AddInt32(&balance, req.ActionNumber)

	log.Println("执行交易：", req.ActionNumber, "，现有余额：", newBalance, "，允许执行")
	b = &proto.Balance{BalanceNumber: newBalance}
	e = nil
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
