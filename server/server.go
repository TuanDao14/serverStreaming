package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/tuanda/serverStreaming/streamingpb"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) Sum(ctx context.Context, req *streamingpb.SumRequest) (*streamingpb.SumResponse, error) {
	log.Println("sum called...")
	resp := &streamingpb.SumResponse{
		Result: req.GetNum1() + req.GetNum2(),
	}

	return resp, nil
}

func (*server) PrimeNumberDecomposition(req *streamingpb.PNDRequest,
	stream streamingpb.CalculatorService_PrimeNumberDecompositionServer) error {
	log.Println("PrimeNumberDecomposition called...")
	k := int32(2)
	N := req.GetNumber()
	for N > 1 {
		if N%k == 0 {
			N = N / k
			//send to client
			stream.Send(&streamingpb.PNDResponse{
				Result: k,
			})
			time.Sleep(1000 * time.Millisecond)
		} else {
			k++
			log.Printf("k increase to %v", k)
		}
	}
	return nil
}

func main() {

	lis, err := net.Listen("tcp", "0.0.0.0:50069")
	if err != nil {
		log.Fatalf("err while create listen %v", err)
	}

	s := grpc.NewServer()

	streamingpb.RegisterCalculatorServiceServer(s, &server{})

	fmt.Println("calculator is running...")
	err = s.Serve(lis)

	if err != nil {
		log.Fatalf("err while serve %v", err)
	}
}
