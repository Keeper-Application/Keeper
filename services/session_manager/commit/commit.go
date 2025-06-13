package commit

import (
	"context"
	"fmt"
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


func (s *SessionManagerServerImpl) CreateSession(ctx context.Context, req *pb.CommitRequest) (*pb.CommitResponse, error) {

	// Check requestor to see empty field. 

	sql := fmt.Sprintf(`
	INSERT INTO sessions(session_id, guardian_id, user_ids, status, session_type)
	VALUES( '%s', '%s' , ARRAY['%s']::UUID[] , '%v' , '%v' )
	`,
	req.SessionInfo.SessionId    ,
	req.SessionInfo.GuardianId   ,
	strings.Join( req.SessionInfo.UserId, "\",") , // Parse the slice next time 
	req.SessionInfo.SessionStatus.String(),
	req.SessionInfo.SessionType.String()  ,
	)

	tag, err := db.Conn.Exec(context.Background() , sql) ; 

	if err != nil {
		fmt.Printf("%v\n", err) ; 
	}

	fmt.Println(sql) ; 
	fmt.Println(tag) ; 

	x := &pb.CommitResponse { 
    CommitStatus  : pb.CommitResponse_S_OK,
    CommitMessage : "it worked",
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


