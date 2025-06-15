package internal

import (
	pgx "github.com/jackc/pgx/v5"
	"log"
	"os"
	ctx "context"
	"github.com/redis/go-redis/v9"
)

var       PSQL_Conn       *pgx.Conn = initializeDB() ; 
var       REDIS_Conn      *redis.Client = initializeRedis() ; 


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

func initializeRedis() *redis.Client {
	opts, err := redis.ParseURL(os.Getenv("REDIS_URL")) ; 
	if err != nil {
		panic(err) ; 
	}
	return redis.NewClient(opts) ; 
}
