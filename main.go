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

// Config holds all configuration variables needed for the application
type Config struct {
	Port        string
	DBURL       string
	Email       string
	EmailSecret string
	TokenSecret string
	RabbitMQURL string
}

// loadConfig loads configuration from environment variables
func loadConfig() (Config, error) {
	if err := godotenv.Load(".env"); err != nil {
		return Config{}, err
	}

	config := Config{
		Port:        os.Getenv("PORT"),
		DBURL:       os.Getenv("DB_URL"),
		Email:       os.Getenv("EMAIL"),
		EmailSecret: os.Getenv("EMAIL_SECRET"),
		TokenSecret: os.Getenv("TOKEN_SECRET"),
		RabbitMQURL: os.Getenv("RABBITMQ_URL"),
	}

	return config, nil
}

// validateConfig ensures all required configuration values are set
func validateConfig(config Config) error {
	if config.Port == "" {
		return log.Output(1, "Set Port in env")
	}

	if config.DBURL == "" {
		return log.Output(1, "Set db connection in env")
	}

	if config.Email == "" {
		return log.Output(1, "Set working email for sending emails")
	}

	if config.EmailSecret == "" {
		return log.Output(1, "Set up Email Secret")
	}

	if config.TokenSecret == "" {
		return log.Output(1, "Set db connection in env")
	}

	if config.RabbitMQURL == "" {
		return log.Output(1, "Set rabbitmq url in env")
	}

	return nil
}

func main() {
	config, err := loadConfig()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	if err := validateConfig(config); err != nil {
		log.Fatal(err)
	}

	lis, err := net.Listen("tcp", config.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	dbConn, err := sql.Open("postgres", config.DBURL)
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}
	dbQueries := database.New(dbConn)
	defer dbConn.Close()

	rabbitmq, err := rabbitmq.NewRabbitMQ(config.RabbitMQURL)
	if err != nil {
		log.Fatal("Can't connect to rabbitmq")
	}

	server := server.NewServer(dbQueries, config.TokenSecret, config.Email, config.EmailSecret, rabbitmq)

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, server)

	reflection.Register(s)
	log.Printf("Server listening on %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
}
