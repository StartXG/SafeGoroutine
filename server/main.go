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
	"time"
)

const PORT = ":8000"

var (
	balance int32 = 1000
	lock    sync.Mutex
	cond    = sync.NewCond(&sync.Mutex{})
)

type server struct {
	proto.UnimplementedBankServiceServer
}

func isTradingHours() bool {
	start := time.Date(0, 0, 0, 17, 0, 0, 0, time.Local)
	end := time.Date(0, 0, 0, 9, 0, 0, 0, time.Local)

	now := time.Now()
	currentTime := time.Date(0, 0, 0, now.Hour(), now.Minute(), now.Second(), now.Nanosecond(), time.Local)

	return currentTime.After(start) && currentTime.Before(end)
}

func (s *server) ModifyNumber(ctx context.Context, req *proto.Action) (b *proto.Balance, e error) {
	cond.L.Lock()
	defer cond.L.Unlock()

	for !isTradingHours() {
		log.Println("非交易时间，禁止交易")
		cond.Wait()
	}

	lock.Lock()
	defer lock.Unlock()

	if balance+req.ActionNumber < 0 {
		log.Println("执行交易：", req.ActionNumber, "，现有余额：", balance, "，余额不足，禁止执行")
		e = fmt.Errorf("余额不足，无法执行操作")
	} else {
		log.Println("执行交易：", req.ActionNumber, "，现有余额：", balance, "，允许执行")
		balance += req.ActionNumber
		e = nil
	}
	b = &proto.Balance{BalanceNumber: balance}
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

	go func() {
		for {
			time.Sleep(1 * time.Minute) // Check every minute
			cond.L.Lock()
			if isTradingHours() {
				cond.Broadcast() // Signal all waiting goroutines if it's trading hours
			}
			cond.L.Unlock()
		}
	}()
}
