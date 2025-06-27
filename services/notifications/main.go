package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/keeper/services/notifications/gen/sessionpb"
	storage "github.com/keeper/services/notifications/internal"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
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

	session_info := &pb.Session{
		SessionStatus  : pb.Session_SESSION_STATUS_ACTIVE,
		SessionType    : pb.Session_SESSION_TYPE_DAILY_RECURRING,
		SessionId      : "826cf6e3-d09a-46f1-9e7c-9d7b3ef3e459",
		GuardianId     : "dfd5910e-942e-491a-a119-3a9ad60d3422",
		UserId         : []string{"02d0f543-b44a-4b88-b8f2-83c1ff5a51ac", "ac884ab5-96ca-482a-9160-56029cfb879f", "884d041c-372d-4ed0-8aeb-f8af49daabfc" },
	}

	sessionreq := &pb.CommitRequest{
    UserType : pb.CommitRequest_GUARDIAN,
    SessionInfo : session_info,
    TenantId        : "",
    RequestorId     : "",
	}

	// Call GetStatus method
	// response, err := client.CreateSession(ctx, sessionreq)
	response, err := client.BeginSession(ctx, sessionreq)
	if err != nil {
		log.Printf("Failed to get lock status: %v", err)
	} else {
		fmt.Printf("Lock status response: %s\n", response.GetCommitStatus())
	}
}

func main() {

	go func() {
		w := storage.KAFKA_Reader ; 
		for {
			m, err := w.ReadMessage(context.Background()); 
			if err != nil {
				/* Write to shared msv buffer */ 
			}
			x := &pb.Session{} ; 
			proto.Unmarshal(m.Value, x) ;
			fmt.Printf("%+v\n\n", x) ;
		}
	}()

  BeginClient("50051") ; 
	for {
		;;
	}

}
