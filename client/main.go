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

func TakeAction(code int) {
	//log.Println("协程", code, "触发")
	//defer log.Println("协程", code, "结束")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
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
	defer conn.Close()
	c = proto.NewBankServiceClient(conn)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go TakeAction(i)
	}
	wg.Wait()
}
