package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/keeper/services/session_manager/sessionpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)



func BeginClient(port string) {
	// Create gRPC client connection
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	// Create LockManager client
	client := pb.NewSessionManagerClient(conn)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create a lock request


	sessionreq := &pb.CommitRequest{
    UserType : pb.CommitRequest_GUARDIAN,
    SessionInfo : &pb.Session{ 
      SessionStatus : pb.Session_SESSION_STATUS_ACTIVE,
      SessionType   : pb.Session_SESSION_TYPE_DAILY_RECURRING,
      SessionId     : "", 
    },
    TenantId        : "",
    RequestorId     : "",
	}

	// Call GetStatus method
	response, err := client.CreateSession(ctx, sessionreq)
	if err != nil {
		log.Printf("Failed to get lock status: %v", err)
	} else {
		fmt.Printf("Lock status response: %s\n", response.GetCommitStatus())
	}
}

func main() {
  fmt.Println("Hello from notifications") ; 
  BeginClient("50051") ; 
}
