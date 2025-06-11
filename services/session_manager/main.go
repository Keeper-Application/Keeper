package main

import (
	"fmt"
	"log"
  "net"

	"google.golang.org/grpc"

	// "keeper/services/session_manager/commit"
	// "keeper/services/lock_manager/issuelock"                   
	pb "github.com/keeper/services/session_manager/gen/sessionpb"                    // For protobuf types
	"github.com/keeper/services/session_manager/commit"                    // For protobuf types
)


func StartServer(port string) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("failed to listen on port %s: %v", port, err)
	}

	grpcServer := grpc.NewServer()
	sessionManagerServer := commit.NewSessionManagerServer() ; 
	
	pb.RegisterSessionManagerServer(grpcServer, sessionManagerServer)
	
	log.Printf("gRPC server listening on port %s", port)
	log.Println("Server starting...")
	
	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve gRPC server: %v", err)
	}

	return nil
}

func main() {
  StartServer("50051") ; 
}
