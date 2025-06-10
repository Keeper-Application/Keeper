package main

import (
	"log"
	"github.com/keeper/services/lock_manager/issuelock"
)

func main() {
	// Call greeting
	
	// Start the gRPC server on port 50051
	if err := issuelock.StartServer("50051"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
