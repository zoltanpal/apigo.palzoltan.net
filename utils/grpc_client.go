package utils

import (
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"golang-restapi/sentimentpb"
)

var (
	sentimentConn   *grpc.ClientConn
	SentimentClient sentimentpb.SentimentServiceClient
)

// InitSentimentClient connects to the gRPC sentiment service and stores the client globally.
func InitSentimentClient(addr string) {
	var err error
	sentimentConn, err = grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to sentiment server at %s: %v", addr, err)
	}
	SentimentClient = sentimentpb.NewSentimentServiceClient(sentimentConn)
}

// CloseSentimentClient closes the gRPC connection (call on shutdown).
func CloseSentimentClient() {
	if sentimentConn != nil {
		_ = sentimentConn.Close()
	}
}
