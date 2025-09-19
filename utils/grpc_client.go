package utils

import (
	"log"

	"google.golang.org/grpc"

	"golang-restapi/sentimentpb"
)

var (
	sentimentConn   *grpc.ClientConn
	SentimentClient sentimentpb.SentimentServiceClient
)

// InitSentimentClient connects to the gRPC sentiment service and stores the client globally.
func InitSentimentClient() {
	var err error
	sentimentConn, err = grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect to sentiment server: %v", err)
	}
	SentimentClient = sentimentpb.NewSentimentServiceClient(sentimentConn)
}

// CloseSentimentClient closes the gRPC connection (call on shutdown).
func CloseSentimentClient() {
	if sentimentConn != nil {
		sentimentConn.Close()
	}
}
