package main

import (
	"fmt"
	"log"
  "net"
	"google.golang.org/grpc"
	pb "github.com/keeper/services/session_manager/gen/sessionpb"
	"github.com/keeper/services/session_manager/commit"         

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
