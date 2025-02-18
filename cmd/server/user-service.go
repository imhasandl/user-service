package server

import (
	"github.com/imhasandl/user-service/internal/database"

	pb "github.com/imhasandl/user-service/protos"
)

type server struct {
	pb.UnimplementedUserServiceServer
	db          *database.Queries
	tokenSecret string
}

func NewServer(dbQueries *database.Queries, tokenSecret string) *server {
	return &server{
		db:          dbQueries,
		tokenSecret: tokenSecret,
	}
}
