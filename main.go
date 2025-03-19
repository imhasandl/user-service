package main

import (
	"database/sql"
	"log"
	"net"
	"os"

	_ "github.com/lib/pq" // Import the postgres driver

	"github.com/imhasandl/user-service/cmd/server"
	"github.com/imhasandl/user-service/internal/database"
	"github.com/imhasandl/user-service/internal/rabbitmq"
	pb "github.com/imhasandl/user-service/protos"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatalf("Set Port in env")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatalf("Set db connection in env")
	}
	
	emailSecret := os.Getenv("EMAIL_SECRET")
	if emailSecret == "" {
		log.Fatalf("Set up Email Secret")
	}

	tokenSecret := os.Getenv("TOKEN_SECRET")
	if tokenSecret == "" {
		log.Fatalf("Set db connection in env")
	}

	rabbitmqURL := os.Getenv("RABBITMQ_URL")
	if rabbitmqURL == "" {
		log.Fatalf("Set rabbitmq url in env")
	}

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listed: %v", err)
	}

	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}
	dbQueries := database.New(dbConn)
	defer dbConn.Close()

	rabbitmq, err := rabbitmq.NewRabbitMQ(rabbitmqURL)
	if err != nil {
		log.Fatal("Can't connect to rabbitmq")
	}

	server := server.NewServer(dbQueries, tokenSecret, rabbitmq)

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, server)

	reflection.Register(s)
	log.Printf("Server listening on %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to lister: %v", err)
	}
}
