package main

import (
	"SafeGoroutine/proto"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

var c proto.BankServiceClient
var wg sync.WaitGroup

func TakeAction() {
	// 为每个协程创建一个上下文，设置超时时间为10秒
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	for i := 0; i < 10; i++ {
		AN := rand.Intn(2000) - 1000
		r, err := c.ModifyNumber(ctx, &proto.Action{ActionNumber: int32(AN)})
		if err != nil {
			log.Println(err.Error())
			continue
		}
		log.Println("执行操作:\t", int32(AN), "\t---> 余额:\t", r.BalanceNumber)
	}
	wg.Done()
}

func main() {
	runtime.GOMAXPROCS(4)
	conn, err := grpc.NewClient("localhost:8000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			panic(err)
		}
	}(conn)
	c = proto.NewBankServiceClient(conn)
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go TakeAction()
	}
	wg.Wait()
}
