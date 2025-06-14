package commit

import (
	"context"
	"fmt"
	"log"

	// "reflect"

	//"log"
	"strings"

	"github.com/jackc/pgx/v5"
	pb "github.com/keeper/services/session_manager/gen/sessionpb"
	db "github.com/keeper/services/session_manager/internal"
)



type         UUID        [16]uint8 ; 

type SessionManagerServerImpl struct {
  pb.UnimplementedSessionManagerServer ; 
}

func NewSessionManagerServer() *SessionManagerServerImpl {
	return &SessionManagerServerImpl{} ; 
}

func buildCreateSessionQuery(req *pb.CommitRequest) string {
	return fmt.Sprintf(`
	INSERT INTO sessions(session_id, guardian_id, user_ids, status, session_type)
	VALUES( '%s', '%s' , ARRAY['%s']::UUID[] , '%v' , '%v' )
	`,
	req.SessionInfo.SessionId    ,
	req.SessionInfo.GuardianId   ,
	strings.Join( req.SessionInfo.UserId, "','") , 
	req.SessionInfo.SessionStatus.String(),
	req.SessionInfo.SessionType.String()  ,
	)
}


// Stub for handling errors later ( Add specialization )

func handleCreateSessionError(req *pb.CommitRequest, e error ) ( *pb.CommitResponse,  error) {
	// Default case for now, add a more descriptive value ( E_UNHANDLED )
	return &pb.CommitResponse{
		CommitStatus: pb.CommitResponse_E_INEXISTENT , 
		CommitMessage: "Unhandled exception occured",
	}, fmt.Errorf("Unhandled exception occured")
}

// TODO: Add error handling 

func (s *SessionManagerServerImpl) CreateSession(ctx context.Context, req *pb.CommitRequest) (*pb.CommitResponse, error) {

	// SQL Query built to hand off to database driver.
	_, err := db.Conn.Exec(context.Background() , buildCreateSessionQuery(req)) ;

	// Call to Error handler function. 
	if err != nil {
		fmt.Errorf("Error occured while creating session: %v", err); 
		return handleCreateSessionError(req, err)  ;
	}

	x := &pb.CommitResponse { 
    CommitStatus  : pb.CommitResponse_S_OK,
    CommitMessage : "Successfully created session",
  }
  return x , nil ; 
}

func handleBeginSessionError(req *pb.CommitRequest, e error ) (*pb.CommitResponse, error) {
	return &pb.CommitResponse{
		CommitStatus: pb.CommitResponse_E_INEXISTENT , 
		CommitMessage: "Unhandled exception occured",
	}, fmt.Errorf("Unhandled exception occured")
}

func buildBeginSessionQuery( req *pb.CommitRequest) string {
	sql := fmt.Sprintf("SELECT * FROM sessions WHERE session_id = '%v'", req.SessionInfo.SessionId ) ; 
	return sql ;
}

func ( s *SessionManagerServerImpl) BeginSession( ctx context.Context, req *pb.CommitRequest) (*pb.CommitResponse, error) {
	// SQL Query built to hand off to database driver.

	fmt.Println(buildBeginSessionQuery(req));
	// res, _ := db.Conn.Query(context.Background(), buildBeginSessionQuery(req)) ; 
	res := db.Conn.QueryRow(context.Background(), buildBeginSessionQuery(req)) ; 

	// res.Next();
	// fmt.Println(res.Values());

	// defer res.Close()
	// Check if result is empty. 
	


	readSessionEntry(&res) ;
	return &pb.CommitResponse{
		CommitStatus:  pb.CommitResponse_S_OK ,
		CommitMessage: "Session successfully began",
	}, nil 
}



func readSessionEntry( row *pgx.Row ) (*pb.Session, error) {

	var session pb.Session ; 

	var session_status string ; 
	var session_type   string  ; 

	// err := (*row).Scan(&session.SessionId , &session.GuardianId, &session.UserId, &session.SessionStatus, &session.SessionType) ;
	
	err := (*row).Scan(&session.SessionId , &session.GuardianId , &session.UserId , &session_status , &session_type) ;
	if err != nil {
		log.Fatalf("Error occured while attempting to deserialize session from db: %v", err) ; 
	}


	session.SessionStatus = pb.Session_SessionStatus(pb.Session_SessionStatus_value[session_status]) ; 
	session.SessionType   = pb.Session_SessionType(pb.Session_SessionType_value[session_type]) ; 

	fmt.Println(session.SessionId) ; 
	fmt.Println(session.UserId) ; 
	fmt.Println(session.GuardianId) ; 
	fmt.Println(session.SessionType) ; 
	fmt.Println(session.SessionStatus) ; 

	return &pb.Session{}, nil ; 
}


// func UpdateSession(req_buffer *CommitRequest , resp_buffer *CommitResponse) CommitStatus {
//   fmt.Println( " Hello from commit " ); 
//   return S_Ok  
// }
//
//
// func DeleteSession(req_buffer *CommitRequest , resp_buffer *CommitResponse) CommitStatus {
//   fmt.Println( " Hello from commit " ); 
//
//   return S_Ok  
// }
//
// func SuspendSession(req_buffer *CommitRequest , resp_buffer *CommitResponse) CommitStatus {
//   fmt.Println( " Hello from commit " ); 
//   return S_Ok  
// }
//
//
//
// func EndSession(req_buffer *CommitRequest , resp_buffer *CommitResponse) CommitStatus {
//   fmt.Println( " Hello from commit " ); 
//   return S_Ok  
// }


