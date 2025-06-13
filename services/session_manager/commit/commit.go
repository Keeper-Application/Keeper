package commit

import (
	"context"
	"fmt"
	"log"
	"strings"

	pb "github.com/keeper/services/session_manager/gen/sessionpb"
	db "github.com/keeper/services/session_manager/internal"
)



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


func handleCreateSessionError(req *pb.CommitRequest, e error ) (*pb.CommitResponse, e error) {
	


	// Default case for now, add a more descriptive value ( E_UNHANDLED )
	return &pb.CommitResponse{
		CommitStatus: pb.CommitResponse_E_INEXISTENT , 
		CommitMessage: "Unhandled exception occured",
	}, fmt.Errorf("Unhandled exception occured")
}

// TODO: Add error handling 

func (s *SessionManagerServerImpl) CreateSession(ctx context.Context, req *pb.CommitRequest) (*pb.CommitResponse, error) {

	sql := buildCreateSessionQuery(req) ;

	_, err := db.Conn.Exec(context.Background() , sql) ;  // Dont really use tag here. 

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
// func BeginSession(req_buffer *CommitRequest , resp_buffer *CommitResponse) CommitStatus {
//   fmt.Println( " Hello from commit " ); 
//   return S_Ok  
// }
//
//
// func EndSession(req_buffer *CommitRequest , resp_buffer *CommitResponse) CommitStatus {
//   fmt.Println( " Hello from commit " ); 
//   return S_Ok  
// }


