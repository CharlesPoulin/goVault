package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"

	"goVault/pkg/kms"
	"goVault/pkg/server"
	"goVault/pkg/storage"

	pb "goVault/proto"
)

func main() {
	// 1. Parse command-line flag
	useKMS := flag.Bool("use-kms", false, "Use AWS KMS for encryption")
	flag.Parse()

	// 2. Initialize storage
	store := storage.NewInMemoryStorage()

	// 3. Initialize the KMS client based on flag
	var kmsClient kms.KMSClient
	if *useKMS {
		fmt.Println("[INFO] Using AWS KMS client...")
		kmsClient = kms.NewAWSKMSClient("arn:aws:kms:us-east-1:123456789012:key/your-master-key-id")
	} else {
		fmt.Println("[INFO] Using local (mock) KMS client...")
		kmsClient = kms.NewLocalKMSClient()
	}

	// 4. Create gRPC server
	grpcServer := grpc.NewServer()

	// 5. Register SecretService
	pb.RegisterSecretServiceServer(grpcServer, server.NewSecretServer(store, kmsClient))

	// 6. Listen on a TCP port
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fmt.Println("Server is running on port :50051")

	// 7. Handle graceful shutdown in a separate goroutine
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh
		fmt.Println("Received shutdown signal... stopping gRPC server gracefully.")
		grpcServer.GracefulStop()
		os.Exit(0)
	}()

	// 8. Serve
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
