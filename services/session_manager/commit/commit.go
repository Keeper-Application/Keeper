package commit

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	pb "github.com/keeper/services/session_manager/gen/sessionpb"
	storage "github.com/keeper/services/session_manager/internal"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
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

/*
   @purpose:      Protobuf Routine which begins a Session with the SessionID Specified by the
									Protobuf Message. 


   @param:        [in]                       *pb.Session   session
                                             Pointer to the session object whose status
                                             is to be validated against the requested action.

   @param:        [in]                       Action        action
                                             The requested operation to apply to the session.

   @param:        [in]                       Action        action
                                             The requested operation to apply to the session.

                                       return

   @return:       nil                      The action is allowed for the session's current state.

   @return:       fmt.Errorf               Descripting string error returned only when the requested 
	                                         action is invalid or inconsistent with the session’s current
																					 status. (e.g., attempting to begin an already active session)

	@notes:         This routine currently performs the following in sequence: Builds
									a database query using the sessionID, Queries the PostgresDB retreiving
									the entry associated with the session. Uses the database entry to construct
									a pb.Session instance, Writes the Contents of the Session into a redis db &
									kafka topic associated with Sessions, and finally returns a CommitResponse
									to the client. TODO: This function should be split up into different parts 
									and be made to run asynchronously. 
*/
func ( s *SessionManagerServerImpl) BeginSession( ctx context.Context, req *pb.CommitRequest) (*pb.CommitResponse, error) {

	// SQL Query built to hand off to database driver.
	res := storage.PSQL_Conn.QueryRow(context.Background(), buildBeginSessionQuery(req)) ; 

	// Create session object & populate it using database response.  
	x := &pb.Session{} ; 

	// TODO: Add error handler function here to remove verbosity. Two cases for now.
	err := sessionFromRow(&res, x) ;

	if err != nil {
		return &pb.CommitResponse{
			CommitStatus:  pb.CommitResponse_E_INEXISTENT ,
			CommitMessage: fmt.Sprintf("Session with ID '%v' does not exist:  %v", req.SessionInfo.SessionId, err) ,
		}, err 
	}

	// Check if Session is in valid state 

	// TODO: Work on this error Handling 
	
	err = reportSession( x , ACTION_BEGIN ) ; 
	if err != nil {
		return &pb.CommitResponse{
			CommitStatus:  pb.CommitResponse_E_INEXISTENT ,
			CommitMessage: fmt.Sprintf("Error occured while beginning session: %v", err) ,
		}, err 
	}

	// Save session to redis db

	err = sessionToRedis( x ) ;

	if err != nil {
		return &pb.CommitResponse{
			CommitStatus:  pb.CommitResponse_E_INEXISTENT ,
			CommitMessage: fmt.Sprintf("Error occured while beginning session: %v", err) ,
		}, fmt.Errorf("Error occured while serializing: %v", err); 
	} 

	err = sessionToKakfa( "my-topic", x ) ;  
	if err != nil {
		return &pb.CommitResponse{
			CommitStatus:  pb.CommitResponse_E_INEXISTENT ,
			CommitMessage: fmt.Sprintf("Error occured while beginning session: %v", err) ,
		} ,fmt.Errorf("Error occured while serializing: %v", err); 
	} 

	return &pb.CommitResponse{
		CommitStatus:  pb.CommitResponse_S_OK ,
		CommitMessage: "Session successfully begun",
	}, nil 

}

func sessionToKakfa( topic string, s *pb.Session) error { 

	// Move writer outisde, it should be configured outside of this function. 

	c := kafka.WriterConfig{
			Brokers: []string{"localhost:9092"},
			Topic: topic,
	}

	w := kafka.NewWriter(c); 

	serialized, err := proto.Marshal(s) ; 
	if err != nil {
		return fmt.Errorf("Error occured while serializing: %v", err); 
	}

	m := kafka.Message{
		Key: []byte(s.SessionId),
		Value: serialized,
	}

	fmt.Printf("%#v", m) ; 
	if err := w.WriteMessages(context.Background(), m); err != nil {
		return fmt.Errorf("Error occured while saving to partition: %v", err); 
	}

	return nil 
}

func sessionToRedis( s *pb.Session ) error {
	// TODO: Marshall this struct once.  && Share between redis & kafka.
	
	serialized, err := proto.Marshal( s ) ;  

	if err != nil {
		return fmt.Errorf("Failed to serialize session data: %v", err); 
	}
	storage.REDIS_Conn.Set( context.Background(), s.SessionId , serialized, 0) ; 
	return nil  ; 
}

func sessionFromRedis(  uuid string ) ( *pb.Session, error ){

	var session pb.Session ; 
	bin, err := storage.REDIS_Conn.Get(context.Background(), uuid).Result() ; 
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve session from RedisDB: %v", err) ; 
	}

	err = proto.Unmarshal([]byte(bin), &session) ; 

	if err != nil {
		return nil, fmt.Errorf("Failed to deserialize session: %v", err) ; 
	}

	return &session, nil ; 
}


func sessionFromRow( row *pgx.Row, session *pb.Session ) (error) {


	var session_status string ; 
	var session_type   string  ; 

	err := (*row).Scan(&session.SessionId , &session.GuardianId , &session.UserId , &session_status , &session_type) ;
	if err != nil {
		return fmt.Errorf("Error occured while de-serializing from db: %v", err) ; 
	}

	session.SessionStatus = pb.Session_SessionStatus(pb.Session_SessionStatus_value[session_status]) ; 
	session.SessionType   = pb.Session_SessionType(pb.Session_SessionType_value[session_type]) ; 
	return nil  ; 
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


