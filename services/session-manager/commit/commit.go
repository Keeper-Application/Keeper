package commit

import (
	"fmt"
	"reflect"
	"time"
	// "unsafe"
)


type                    Opcode                   uint16 ; 

type                    RoleHint                 uint8  ; 

type                    SessionState             uint16 ; 

type                    SessionType              uint16 ; 

type                    CommitStatus             uint16 ; 

type                    CommitHandler            func ( *CommitRequest , *CommitResponse ) CommitStatus ; 


type     CommitRequest      struct {
    Hint              RoleHint           ;   // 1
    Opcode            Opcode             ;   // 2
    Payload           *byte              ;   // 8
    TenantID          [16]byte           ;   // 16
    SessionID         [16]byte           ;   // 16
    RequestorID       [16]byte           ;   // 16 
}; 

type     CommitResponse     struct {
    Status            int                ;
    SessionState      SessionState       ;
    Message           string             ;
    Output            map[string]any     ;
    Timestamp         time.Time          ;
}


const (
	SessionTypeUndefined SessionType   =          iota  // fallback/default
	SessionTypeTimeBound                                // e.g. 3pm–5pm, recurring or one-time
	SessionTypeGoalBased                                // e.g. complete a checklist, task, or app usage limit
	SessionTypeAppSpecific                              // block apps until specific apps are used (e.g. Duolingo)
	SessionTypeScreenTimeBudget                         // session ends when a time budget is exhausted (e.g. 60 minutes of YouTube)
	SessionTypeLocationBased                            // triggered only in specific geofenced areas (school, library)
	SessionTypeDailyRecurring                           // scheduled daily lock window (e.g. every day 10pm–7am)
	SessionTypeOneTimeChallenge                         // Guardian sets a single-use challenge (e.g., "1 hour writing")
)

const (
  OpcodeUndefined       Opcode        =          iota 
  OpcodeCreateSession
  OpcodeUpdateSession
  OpcodeDeleteSession
  OpcodeSuspendSession
  OpcodeEndSession 
  OpcodeBeginSession 
)

const (
  SessionStatusErr      SessionState   =          iota  
  SessionStatusActive  
  SessionStatusEnded 
  SessionSuspended
  SessionExpired  
)


const (
  E_Exists              CommitStatus   =          iota  
  E_Access         
  E_LimitReached
  E_Inexistent 
  E_Permission 
  E_Busy
  S_Ok 
)


var dispatch_map = map[Opcode]CommitHandler{
  OpcodeCreateSession  : CreateSession ,
  OpcodeUpdateSession  : UpdateSession ,
  OpcodeBeginSession   : BeginSession  ,
  OpcodeEndSession     : BeginSession  ,
  OpcodeDeleteSession  : DeleteSession , 
  OpcodeSuspendSession : SuspendSession,
}

/*
   @purpose:      Creates an instance of a session & saves it to the database 

   @param:        [in]                     CommitResquest*    req_buffer

   @param:        [in]                     CommitResponse*    resp_buffer

                                       return

   @code:         E_Exists                 Session already exists. 


   @code:         E_LimitReached           Requester has reached limits of the amount of 
                                           sessions ( Four for standard users ) they are 
                                           allowed to be an orchestrator of. 
  

   @code:         S_Ok                     Session created successfully. 

   @notes:        No notes for now.  
*/


func CreateSession(req_buffer *CommitRequest , resp_buffer *CommitResponse) CommitStatus {
  var s  CommitRequest ; 
  fmt.Println(reflect.TypeOf(s).Align()) ; 
  fmt.Println(reflect.TypeOf(s).Size()) ; 
  return S_Ok  
}


func UpdateSession(req_buffer *CommitRequest , resp_buffer *CommitResponse) CommitStatus {
  fmt.Println( " Hello from commit " ); 
  return S_Ok  
}


func DeleteSession(req_buffer *CommitRequest , resp_buffer *CommitResponse) CommitStatus {
  fmt.Println( " Hello from commit " ); 

  return S_Ok  
}

func SuspendSession(req_buffer *CommitRequest , resp_buffer *CommitResponse) CommitStatus {
  fmt.Println( " Hello from commit " ); 
  return S_Ok  
}

func BeginSession(req_buffer *CommitRequest , resp_buffer *CommitResponse) CommitStatus {
  fmt.Println( " Hello from commit " ); 
  return S_Ok  
}


func EndSession(req_buffer *CommitRequest , resp_buffer *CommitResponse) CommitStatus {
  fmt.Println( " Hello from commit " ); 
  return S_Ok  
}

// Greedy approach 
// Check size of previous character and alignment of self 
// alignment = align 
// everytime we see a type we subtract from the align and go to the next 
// If the value of the next type would cause us to exceed our alignment add padding and repeat. 

// Suppose Align = 4 & F1 = uint8 that means we subtract 1 and go forward
// now two cases arise, whether continuing would cause us to overstep our alignment bounds ( given by F + 1 ) case we would have to add padding.
// or we can continue 

func ReprPadding( value interface{} ) { // Why shtty align of non-power of 2. keep it @ 2^x
  t := reflect.TypeOf(value) ; 
  for i := t.NumField() - 1 ; i >= 0 ; i-- {
    // if t.Field(i + 1 ).Offset > t.Field(i).Type.Size() {
    //   fmt.Print("Padding was added")
    // }

    fmt.Println(t.Field(i).Type.Size()) ;
  }
  fmt.Printf("Size: %v\n", reflect.TypeOf(value).Size()) ;  
}

func alignment_as_int(uip uintptr) int {
    return int(uip)
}
