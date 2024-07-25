package main

import (
	"SafeGoroutine/proto"
	"context"
	"log"
	"math/rand"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var c proto.BankServiceClient
var wg sync.WaitGroup
var mu sync.Mutex

const workerNumbers = 2

func TakeAction(workerID int, job int) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	AN := rand.Intn(2000) - 1000
	r, err := c.ModifyNumber(ctx, &proto.Action{ActionNumber: int32(AN)})
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Println("workerID:", workerID, "job:", job, "执行操作:\t", int32(AN), "\t---> 余额:\t", r.BalanceNumber)
}

func main() {
	var executed int = 0
	conn, err := grpc.NewClient("localhost:8000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	c = proto.NewBankServiceClient(conn)
	jobs := make(chan int, 2)

	for w := 1; w <= workerNumbers; w++ {
		wg.Add(1)
		go func(workerID int, jobs <-chan int) {
			defer wg.Done()
			for job := range jobs {
				TakeAction(workerID, job)

				mu.Lock()
				executed += 1
				mu.Unlock()
			}
		}(w, jobs)
	}

	log.Println("Sending jobs...")
	for i := 0; i < 100; i++ {
		jobs <- i
	}

	close(jobs)
	log.Println("All jobs sent.")

	wg.Wait()
	log.Println("All workers done.", executed)

}
