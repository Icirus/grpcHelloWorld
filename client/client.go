package main

import (
	"fmt"
	"grpc_hello_world/chat"
	"log"
	"math/rand"
	"sync"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
    b := make([]byte, n)
    for i := range b {
        b[i] = letterBytes[rand.Intn(len(letterBytes))]
    }
    return string(b)
}


func generateMessage (lenOfMessage int, messageCount int) (message *chat.Message) {
	// function to generate the randomized message
	message = &chat.Message{
		Body: RandStringBytes(lenOfMessage),
		MessageNumber: int64(messageCount),
		Timestamps: time.Now().Unix(),
	}
	return message
}

func calculateLatency (sendTime int64) (int64){
		
		return time.Now().Unix()-sendTime
}

func callSayHello(conn *grpc.ClientConn, messageCount int, wg *sync.WaitGroup, sem chan int){
	defer wg.Done()
	c := chat.NewChatServiceClient(conn)
	response, err := c.SayHello(context.Background(), generateMessage(25, messageCount))
	if err != nil{
		log.Fatalf("Error when calling SayHello: %s", err)
	}

	log.Printf("Response From Server: %v -- Roundtrip time: %d", response, calculateLatency(response.Timestamps))
	fmt.Println(sem)
	<-sem
}

func main() {
	var conn *grpc.ClientConn
	var wg sync.WaitGroup
	var interationCount int
	const MAX = 2000
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %s", err)
	}
	defer conn.Close()
	interationCount = 1000000
	wg.Add(interationCount)

	messagesSent := 0
	sem := make(chan int, MAX)
	for messagesSent < interationCount {
		sem <- 1
		go callSayHello(conn, messagesSent, &wg, sem)
		messagesSent += 1
	}

	wg.Wait()
}