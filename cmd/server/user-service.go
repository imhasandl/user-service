package server

import (
	"context"

	"github.com/google/uuid"
	postService "github.com/imhasandl/post-service/cmd/helper"
	"github.com/imhasandl/user-service/internal/database"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	authService "github.com/imhasandl/auth-service/cmd/helper/auth"

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

func (s *server) GetUserByEmailOrUsername(
	ctx context.Context,
	req *pb.GetUserByEmailOrUsernameRequest,
) (*pb.GetUserByEmailOrUsernameResponse, error) {
	userParams := database.GetUserByEmailOrUsernameParams{
		Email:    req.GetIdentifier(),
		Username: req.GetIdentifier(),
	}

	user, err := s.db.GetUserByEmailOrUsername(ctx, userParams)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "can't get user from db: %v - GetUserByEmailOrUsername", err)
	}

	return &pb.GetUserByEmailOrUsernameResponse{
		User: &pb.User{
			Id:        user.ID.String(),
			CreatedAt: timestamppb.New(user.CreatedAt),
			UpdatedAt: timestamppb.New(user.UpdatedAt),
			Email:     user.Email,
			Username:  user.Username,
			IsPremium: user.IsPremium,
		},
	}, nil
}

func (s *server) GetUserByID(ctx context.Context, req *pb.GetUserByIDRequest) (*pb.GetUserByIDResponse, error) {
	userID, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "can't parse user id from incomin request: %v - GetUserByID", err)
	}

	user, err := s.db.GetUserById(ctx, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "can't get user from db: %v - GetUserByID", err)
	}

	return &pb.GetUserByIDResponse{
		User: &pb.User{
			Id: user.ID.String(),
			CreatedAt: timestamppb.New(user.CreatedAt),
			UpdatedAt: timestamppb.New(user.UpdatedAt),
			Email: user.Email,
			Username: user.Username,
			IsPremium: user.IsPremium,
		},
	}, nil
}
func (s *server) ChangeUsername(ctx context.Context, req *pb.ChangeUsernameRequest) (*pb.ChangeUsernameResponse, error) {
	accessToken, err := postService.GetBearerTokenFromGrpc(ctx)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "can't get token from header: %v - ChangeUsername", err)
	}

	userID, err := postService.ValidateJWT(accessToken, s.tokenSecret)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "you are not allowed to change username: %v - ChangeUsername", err)
	}

	changeUsernameParams := database.ChangeUsernameParams{
		ID:       userID,
		Username: req.GetUsername(),
	}

	user, err := s.db.ChangeUsername(ctx, changeUsernameParams)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "can't change username in db: %v - ChangeUsername", err)
	}

	return &pb.ChangeUsernameResponse{
		User: &pb.User{
			Id:        user.ID.String(),
			CreatedAt: timestamppb.New(user.CreatedAt),
			UpdatedAt: timestamppb.New(user.UpdatedAt),
			Email:     user.Email,
			Username:  user.Username,
			IsPremium: user.IsPremium,
		},
	}, nil
}

func (s *server) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {
	accessToken, err := postService.GetBearerTokenFromGrpc(ctx)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "can't get authorization token from header: %v - ChangePassword", err)
	}

	userID, err := postService.ValidateJWT(accessToken, s.tokenSecret)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "you are not allowed to delete user: %v - ChangePassword", err)
	}	

	hashedPassword, err := authService.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "can't hash password: %v - ChangePassword", err)
	}

	changePasswordParams := database.ChangePasswordParams{
		ID: userID,
		Password: hashedPassword,
	}

	err = s.db.ChangePassword(ctx, changePasswordParams)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "can't change password: %v - ChangePassword", err)
	}

	return &pb.ChangePasswordResponse{
		Status: "Password changed successfully",
	}, nil
}

func (s *server) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	accessToken, err := postService.GetBearerTokenFromGrpc(ctx)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "can't get authorization token from header: %v - DeleteUser", err)
	}

	userID, err := postService.ValidateJWT(accessToken, s.tokenSecret)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "you are not allowed to delete user: %v - DeleteUser", err)
	}

	_, err = s.db.DeleteUser(ctx, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "can't delete user from db: %v - DeleteUser", err)
	}

	return &pb.DeleteUserResponse{
		Status: "user delete successfully",
	}, nil
}