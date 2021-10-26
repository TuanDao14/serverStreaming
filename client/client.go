package main

import (
	"context"
	"io"
	"log"

	"github.com/tuanda/serverStreaming/streamingpb"
	"google.golang.org/grpc"
)

func main() {
	cc, err := grpc.Dial("localhost:50069", grpc.WithInsecure())

	if err != nil {
		log.Fatalf(" err while dial %v", err)
	}
	defer cc.Close()

	client := streamingpb.NewCalculatorServiceClient(cc)

	log.Printf("service client %f\n", client)
	callSum(client)
	callPND(client)
}

func callSum(c streamingpb.CalculatorServiceClient) {
	log.Println("calling sum api")
	resp, err := c.Sum(context.Background(), &streamingpb.SumRequest{
		Num1: 7,
		Num2: 6,
	})

	if err != nil {
		log.Fatalf("call sum api err %v", err)
	}

	log.Printf("sum api response %v\n", resp.GetResult())
}

func callPND(c streamingpb.CalculatorServiceClient) {
	log.Println("calling pnd api")
	stream, err := c.PrimeNumberDecomposition(context.Background(), &streamingpb.PNDRequest{
		Number: 120,
	})

	if err != nil {
		log.Fatalf("callPND err %v", err)
	}

	for {
		resp, recvErr := stream.Recv()
		if recvErr == io.EOF {
			log.Println("server finish streaming")
			return
		}

		log.Printf("prime number %v", resp.GetResult())
	}
}
