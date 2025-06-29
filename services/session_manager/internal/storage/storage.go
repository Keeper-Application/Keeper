package internal

import (
	pgx "github.com/jackc/pgx/v5"
	"log"
	"os"
	ctx "context"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
)


var       PSQL_Conn       *pgx.Conn      = initializeDB() ; 
var       REDIS_Conn      *redis.Client  = initializeRedis() ; 
var       KAFKA_Writer    *kafka.Writer  = initializeKafkaWriter() ; 
var       KAFKA_Reader    *kafka.Reader  = initializeKafkaReader() ; 



type    QueryBuilder   interface {
	ToQuery(string) string ; 
}

/*
   @purpose:      Returns an instance of a PostgreSQL Client connection to be used
	 								throughout the lifetime of the server. It utilizes the hosts environment 
									variables to configure the connection. 

                                       return

   @return:       *pgx.Conn                A pointer to a PostgreSQL client connection. 


	@notes:         This routine is used to initializes a database connection to be used throughout the 
									session. Curretnly just uses the host env to configure but later on should be moved 
									to use a .env file. 
*/

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


/*
   @purpose:      Returns an instance of a Redis Client connection to be used
	 								throughout the lifetime of the server. It utilizes the hosts environment 
									variables to configure the connection. 

                                       return

   @return:       *redis.client       Pointer to a redis client connection.


	@notes:         This routine is used to initializes a redis connection to be used throughout the 
									session. Currently just uses the host env to configure but later on should be moved 
									to use a .env file. 
*/

func initializeRedis() *redis.Client {
	if len(os.Getenv("REDIS_URL")) == 0 {
		log.Fatal("Malformed environment variable: $REDIS_URL is not defined.") ; 
	}
	opts, err := redis.ParseURL(os.Getenv("REDIS_URL")) ; 
	if err != nil {
		panic(err) ; 
	}
	return redis.NewClient(opts) ; 
}

/*
   @purpose:      Returns an instance of a KafkaWriter Client connection to be used
	 								throughout the lifetime of the server. It utilizes the hosts environment 
									variables to configure the connection. 

                                       return

   @return:       *kafka.Writer       The action is allowed for the session's current state.


	@notes:         This routine is used to initializes a Kafka connection to be used throughout the 
									session. Currrently just uses the host env to configure but later on should be moved 
									to use a .env file. 
*/
func initializeKafkaWriter() *kafka.Writer{
	c := kafka.WriterConfig{
		Brokers: []string{"localhost:9092"},
		Topic: "my-topic",
	}
	return kafka.NewWriter(c) ;
}

/*
   @purpose:      Returns an instance of a Kafka Client connection to be used throughout the 
	 								throughout the lifetime of the server. It utilizes the hosts environment 
									variables to configure the connection. 

                                       return

   @return:       *kafka.Reader        The action is allowed for the session's current state.


	@notes:         This routine is used to initializes a database connection to be used throughout the 
									session. Currently just uses the host env to configure but later on should be moved 
									to use a .env file. 
*/

func initializeKafkaReader() *kafka.Reader{
	c := kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Partition: 0, 
		Topic: "my-topic",
		// MaxBytes: 1 << 20 ,   // Uncomment to set MaxBytes to read ( Truncates whats left of a message > 1MB ).
	}
	return kafka.NewReader(c) ;
}
