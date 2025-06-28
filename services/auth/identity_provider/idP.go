package identityprovider

import (
	"context"
	"fmt"

	storage "github.com/keeper/services/auth/internal"
)

type UUID string ;


// Empty QueryBuilder interface for now. 

type QueryBuilder interface {
	~string  ;
}

// For use later, can make this a protobuf datastruct.  This serves to store the users retrieved from the db.

type User struct{
	user_id UUID ; 
} ; 

func fetchUser( id UUID ) () { 
	x := User{ user_id: "02d0f543-b44a-4b88-b8f2-83c1ff5a51ac" } ; 
	row :=  storage.PSQL_Conn.QueryRow( context.Background(), fmt.Sprintf("SELECT * FROM Users WHERE user_id=%s", id));
	row.Scan(&x.user_id) ;
	fmt.Println("User_id: %v", x.user_id) ;
}

func AuthorizeUser( id UUID ) { 
	fetchUser(id) ;
}
