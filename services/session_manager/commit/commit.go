package commit

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	pb "github.com/keeper/services/session_manager/gen/sessionpb"
	storage "github.com/keeper/services/session_manager/internal"
)



type         UUID        [16]uint8 ; 

type         Action      uint8     ; 

const (
	ACTION_UNDEFINED   Action  = iota
	ACTION_CREATE
	ACTION_UPDATE
	ACTION_BEGIN
	ACTION_END
	ACTION_DELETE
	ACTION_SUSPEND
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
	_, err := storage.PSQL_Conn.Exec(context.Background() , buildCreateSessionQuery(req)) ;

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

/*
   @purpose:      Validates whether a session can safely undergo the given action.
                  Ensures business logic constraints (e.g., no reactivating an active session)
                  are respected before proceeding with session lifecycle operations.

   @typedef:      Action -                   Enum-like custom type representing session lifecycle 
	                                           operations such as create, update, begin, end, delete,
																						 suspend.

   @param:        [in]                       *pb.Session   session
                                             Pointer to the session object whose status
                                             is to be validated against the requested action.

   @param:        [in]                       Action        action
                                             The requested operation to apply to the session.

                                       return

   @code:         nil                      The action is allowed for the session's current state.

   @code:         fmt.Errorf               Descripting string error returned only when the requested 
	                                         action is invalid or inconsistent with the session’s current
																					 status. (e.g., attempting to begin an already active session)

   @notes:        This function currently enforces validation only for ACTION_BEGIN.
                  All other actions are treated as valid operations by default.
                  Additional validation logic may be implemented in the future for stricter
                  session state management.
*/

// WARNING: This introduces statefulness into the code. 
func reportSession(session *pb.Session , action Action ) error {
	switch action {

	case ACTION_CREATE :
		return nil 
	case ACTION_UPDATE :
		return nil
	case ACTION_BEGIN :
		if session.SessionStatus == pb.Session_SESSION_STATUS_ACTIVE{

			return fmt.Errorf("Attempted to begin a session which is already active")
		}
		return nil
	case ACTION_END :
		if session.SessionStatus == pb.Session_SESSION_SUSPENDED || 
		   session.SessionStatus == pb.Session_SESSION_EXPIRED   || 
			 session.SessionStatus == pb.Session_SESSION_STATUS_ENDED {
			return fmt.Errorf("Attempted to delete an active session")
		}
		return nil
	case ACTION_DELETE :
		if session.SessionStatus == pb.Session_SESSION_STATUS_ACTIVE || 
		   session.SessionStatus == pb.Session_SESSION_SUSPENDED     {
			return fmt.Errorf("Attempted to delete an active/suspended session")
		}
		return nil
	case ACTION_SUSPEND :
		return nil
	default:
		return fmt.Errorf("Undefined action")
	}
}

func ( s *SessionManagerServerImpl) BeginSession( ctx context.Context, req *pb.CommitRequest) (*pb.CommitResponse, error) {

	// SQL Query built to hand off to database driver.
	res := storage.PSQL_Conn.QueryRow(context.Background(), buildBeginSessionQuery(req)) ; 


	// Create session object & populate it using database response.  
	x := &pb.Session{} ; 



	// Handle error if it occured. (Have to add a better default value for enums mane)
	// TODO: Add error handler function here to remove verbosity. Two cases for now.
	
	err := sessionFromRow(&res, x) ;
	if err != nil {
		switch err {
		case pgx.ErrNoRows:
		return &pb.CommitResponse{
			CommitStatus:  pb.CommitResponse_E_INEXISTENT ,
			CommitMessage: fmt.Sprintf("Session with ID '%v' does not exist:  %v", req.SessionInfo.SessionId, err) ,
		}, err 
		default:
		return &pb.CommitResponse{
			CommitStatus:  pb.CommitResponse_E_INEXISTENT ,
			CommitMessage: fmt.Sprintf("Error occured while beginning session: %v", err) ,
		}, err 
		}
	}

	// Check if Session is in valid state 

	err = reportSession( x , ACTION_BEGIN ) ; 
	if err != nil {
		return &pb.CommitResponse{
			CommitStatus:  pb.CommitResponse_E_INEXISTENT ,
			CommitMessage: fmt.Sprintf("Error occured while beginning session: %v", err) ,
		}, err 
	}

	// Save session to redis db
	



	// Create log within kafka topic ( for LockManager & Notifications microservices ) 


	// Update state within db. 
	
	return &pb.CommitResponse{
		CommitStatus:  pb.CommitResponse_S_OK ,
		CommitMessage: "Session successfully begun",
	}, nil 

}


func sessionFromRow( row *pgx.Row, buffer *pb.Session ) (error) {

	var session pb.Session ; 

	var session_status string ; 
	var session_type   string  ; 

	err := (*row).Scan(&session.SessionId , &session.GuardianId , &session.UserId , &session_status , &session_type) ;
	if err != nil {
		return fmt.Errorf("Error occured while de-serializing from db: %v", err) ; 
	}

	session.SessionStatus = pb.Session_SessionStatus(pb.Session_SessionStatus_value[session_status]) ; 
	session.SessionType   = pb.Session_SessionType(pb.Session_SessionType_value[session_type]) ; 

	fmt.Println(session.SessionId) ; 
	fmt.Println(session.UserId) ; 
	fmt.Println(session.GuardianId) ; 
	fmt.Println(session.SessionType) ; 
	fmt.Println(session.SessionStatus) ; 

	return nil ; 
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


