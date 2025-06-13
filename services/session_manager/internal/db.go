package database 

import (
	pgx "github.com/jackc/pgx/v5"
	"log"
	"os"
	ctx "context"
)

var Conn *pgx.Conn = initializeDB() ; 


func initializeDB() *pgx.Conn {
	if len(os.Getenv("DB_URL")) == 0 {
		log.Fatal("Malformed environment variable: $DB_URL is not defined.") ; 
	}
	conn, err := pgx.Connect( ctx.Background(), os.Getenv("DB_URL"))
	if err != nil {
		log.Fatalf("Error occured while establishing connection: %v",  err) ; 
	}

	return conn ; 
}

