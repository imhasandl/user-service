package server

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	authService "github.com/imhasandl/auth-service/cmd/auth"
	postService "github.com/imhasandl/post-service/cmd/auth"
	"github.com/imhasandl/user-service/internal/database"
	"github.com/imhasandl/user-service/internal/rabbitmq"
	"github.com/streadway/amqp"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/timestamppb"

	helper "github.com/imhasandl/user-service/cmd/helper"
	pb "github.com/imhasandl/user-service/protos"
)

// UserServer defines the interface for the user service server
type UserServer interface {
	pb.UserServiceServer
}

type server struct {
	pb.UnimplementedUserServiceServer
	db          *database.Queries
	tokenSecret string
	email       string
	emailSecret string
	rabbitmq    *rabbitmq.RabbitMQ
}

// NewServer creates and returns a new instance of the user service server.
// It initializes the server with the provided repository, config, and optional handler.
func NewServer(dbQueries *database.Queries, tokenSecret string, email string, emailSecret string, rabbitmq *rabbitmq.RabbitMQ) UserServer {
	return &server{
		db:          dbQueries,
		tokenSecret: tokenSecret,
		email:       email,
		emailSecret: emailSecret,
		rabbitmq:    rabbitmq,
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
		return nil, helper.RespondWithErrorGRPC(ctx, codes.Internal, "can't get user from db: GetUserByEmailOrUsername", err)
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
		return nil, helper.RespondWithErrorGRPC(ctx, codes.Internal, "can't parse user id from incoming request: GetUserByID", err)
	}

	user, err := s.db.GetUserById(ctx, userID)
	if err != nil {
		return nil, helper.RespondWithErrorGRPC(ctx, codes.Internal, "can't get user from db: GetUserByID", err)
	}

	return &pb.GetUserByIDResponse{
		User: &pb.User{
			Id:               user.ID.String(),
			CreatedAt:        timestamppb.New(user.CreatedAt),
			UpdatedAt:        timestamppb.New(user.UpdatedAt),
			Email:            user.Email,
			Username:         user.Username,
			IsPremium:        user.IsPremium,
			VerificationCode: user.VerificationCode,
			IsVerified:       user.IsVerified,
		},
	}, nil
}

func (s *server) GetUserByToken(ctx context.Context, req *pb.GetUserByTokenRequest) (*pb.GetUserByTokenResponse, error) {
	accessToken, err := postService.GetBearerTokenFromGrpc(ctx)
	if err != nil {
		return nil, helper.RespondWithErrorGRPC(ctx, codes.InvalidArgument, "can't get token from header: GetUserByToken", err)
	}

	userID, err := postService.ValidateJWT(accessToken, s.tokenSecret)
	if err != nil {
		return nil, helper.RespondWithErrorGRPC(ctx, codes.Unauthenticated, "can't validate token: GetUserByToken", err)
	}

	user, err := s.db.GetUserById(ctx, userID)
	if err != nil {
		return nil, helper.RespondWithErrorGRPC(ctx, codes.Internal, "can't get user from db: GetUserByToken", err)
	}

	return &pb.GetUserByTokenResponse{
		User: &pb.User{
			Id:               user.ID.String(),
			CreatedAt:        timestamppb.New(user.CreatedAt),
			UpdatedAt:        timestamppb.New(user.UpdatedAt),
			Email:            user.Email,
			Username:         user.Username,
			IsPremium:        user.IsPremium,
			VerificationCode: user.VerificationCode,
			IsVerified:       user.IsVerified,
		},
	}, nil
}

func (s *server) GetAllUsers(ctx context.Context, req *pb.GetAllUsersRequest) (*pb.GetAllUsersResponse, error) {
	users, err := s.db.GetAllUsers(ctx)
	if err != nil {
		return nil, helper.RespondWithErrorGRPC(ctx, codes.Internal, "can't get users from db: GetAllUsers", err)
	}

	pbUsers := make([]*pb.User, len(users))
	for i, user := range users {
		pbUsers[i] = &pb.User{
			Id:               user.ID.String(),
			CreatedAt:        timestamppb.New(user.CreatedAt),
			UpdatedAt:        timestamppb.New(user.UpdatedAt),
			Email:            user.Email,
			Username:         user.Username,
			IsPremium:        user.IsPremium,
			VerificationCode: user.VerificationCode,
			IsVerified:       user.IsVerified,
		}
	}

	return &pb.GetAllUsersResponse{
		Users: pbUsers,
	}, nil
}

func (s *server) ChangeUsername(ctx context.Context, req *pb.ChangeUsernameRequest) (*pb.ChangeUsernameResponse, error) {
	accessToken, err := postService.GetBearerTokenFromGrpc(ctx)
	if err != nil {
		return nil, helper.RespondWithErrorGRPC(ctx, codes.InvalidArgument, "can't get token from header: ChangeUsername", err)
	}

	userID, err := postService.ValidateJWT(accessToken, s.tokenSecret)
	if err != nil {
		return nil, helper.RespondWithErrorGRPC(ctx, codes.Unauthenticated, "you are not allowed to change username: ChangeUsername", err)
	}

	changeUsernameParams := database.ChangeUsernameParams{
		ID:       userID,
		Username: req.GetUsername(),
	}

	user, err := s.db.ChangeUsername(ctx, changeUsernameParams)
	if err != nil {
		return nil, helper.RespondWithErrorGRPC(ctx, codes.Internal, "can't change username in db: ChangeUsername", err)
	}

	return &pb.ChangeUsernameResponse{
		User: &pb.User{
			Id:               user.ID.String(),
			CreatedAt:        timestamppb.New(user.CreatedAt),
			UpdatedAt:        timestamppb.New(user.UpdatedAt),
			Email:            user.Email,
			Username:         user.Username,
			IsPremium:        user.IsPremium,
			VerificationCode: user.VerificationCode,
			IsVerified:       user.IsVerified,
		},
	}, nil
}

func (s *server) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {
	accessToken, err := postService.GetBearerTokenFromGrpc(ctx)
	if err != nil {
		return nil, helper.RespondWithErrorGRPC(ctx, codes.InvalidArgument, "can't get authorization token from header: ChangePassword", err)
	}

	userID, err := postService.ValidateJWT(accessToken, s.tokenSecret)
	if err != nil {
		return nil, helper.RespondWithErrorGRPC(ctx, codes.Unauthenticated, "you are not allowed to delete user: ChangePassword", err)
	}

	hashedPassword, err := authService.HashPassword(req.GetPassword())
	if err != nil {
		return nil, helper.RespondWithErrorGRPC(ctx, codes.Internal, "can't hash password: ChangePassword", err)
	}

	changePasswordParams := database.ChangePasswordParams{
		ID:       userID,
		Password: hashedPassword,
	}

	err = s.db.ChangePassword(ctx, changePasswordParams)
	if err != nil {
		return nil, helper.RespondWithErrorGRPC(ctx, codes.Internal, "can't change password: ChangePassword", err)
	}

	return &pb.ChangePasswordResponse{
		Status: "Password changed successfully",
	}, nil
}

func (s *server) SubscribeUser(ctx context.Context, req *pb.SubscribeUserRequest) (*pb.SubscribeUserResponse, error) {
	accessToken, err := postService.GetBearerTokenFromGrpc(ctx)
	if err != nil {
		return nil, helper.RespondWithErrorGRPC(ctx, codes.InvalidArgument, "can't token from ctx - SubscribeUser", err)
	}

	subscriberUserID, err := postService.ValidateJWT(accessToken, s.tokenSecret)
	if err != nil {
		return nil, helper.RespondWithErrorGRPC(ctx, codes.InvalidArgument, "can't validate provided token - SubscribeUser", err)
	}

	subscribedUserID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, helper.RespondWithErrorGRPC(ctx, codes.InvalidArgument, "can't parse user's uuid - SubscribeUser", err)
	}

	subscribeUserParamsParams := database.SubscribeUserParams{
		ID:          subscribedUserID,
		ArrayAppend: subscriberUserID,
	}

	err = s.db.SubscribeUser(ctx, subscribeUserParamsParams)
	if err != nil {
		return nil, helper.RespondWithErrorGRPC(ctx, codes.Internal, "can't sub to user - SubscribeUser", err)
	}

	messageJSON, err := json.Marshal(map[string]string{
		"title":           "New Notification",
		"sender_username": subscribedUserID.String(),
		"receiver_id":     subscribedUserID.String(),
		"content":         fmt.Sprintf("This user %v subscribed on you", subscriberUserID),
		"sent_at":         time.Now().GoString(),
	})
	if err != nil {
		return nil, helper.RespondWithErrorGRPC(ctx, codes.Internal, "can't send subscribe notification - SubscribeUser", err)
	}

	err = s.rabbitmq.Channel.Publish(
		rabbitmq.ExchangeName, // exchange
		rabbitmq.RoutingKey,   // routing key
		false,                 // mandatory
		false,                 // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        messageJSON,
		})
	if err != nil {
		return nil, helper.RespondWithErrorGRPC(ctx, codes.Internal, "can't publish message to RabbitMQ - SubscribeUser", err)
	}

	return &pb.SubscribeUserResponse{
		Status: true,
	}, nil
}

func (s *server) UnsubscribeUser(ctx context.Context, req *pb.UnsubscribeUserRequest) (*pb.UnsubscribeUserReponse, error) {
	accessToken, err := postService.GetBearerTokenFromGrpc(ctx)
	if err != nil {
		return nil, helper.RespondWithErrorGRPC(ctx, codes.InvalidArgument, "can't token from ctx - UnsubscribeUser", err)
	}

	unSubscriberUserID, err := postService.ValidateJWT(accessToken, s.tokenSecret)
	if err != nil {
		return nil, helper.RespondWithErrorGRPC(ctx, codes.InvalidArgument, "can't validate provided token - UnsubscribeUser", err)
	}

	unSubscribedUserID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, helper.RespondWithErrorGRPC(ctx, codes.InvalidArgument, "can't parse user's uuid - UnsubscribeUser", err)
	}

	unSubscribeUserParamsParams := database.UnsubscribeUserParams{
		ID:          unSubscribedUserID,
		ArrayRemove: unSubscriberUserID,
	}

	err = s.db.UnsubscribeUser(ctx, unSubscribeUserParamsParams)
	if err != nil {
		return nil, helper.RespondWithErrorGRPC(ctx, codes.Internal, "can't unsub to user - UnsubscribeUser", err)
	}

	return &pb.UnsubscribeUserReponse{
		Status: true,
	}, nil
}

func (s *server) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	accessToken, err := postService.GetBearerTokenFromGrpc(ctx)
	if err != nil {
		return nil, helper.RespondWithErrorGRPC(ctx, codes.InvalidArgument, "can't get authorization token from header: DeleteUser", err)
	}

	userID, err := postService.ValidateJWT(accessToken, s.tokenSecret)
	if err != nil {
		return nil, helper.RespondWithErrorGRPC(ctx, codes.Unauthenticated, "you are not allowed to delete user: DeleteUser", err)
	}

	user, err := s.db.GetUserById(ctx, userID)
	if err != nil {
		return nil, helper.RespondWithErrorGRPC(ctx, codes.Internal, "can't get user from db: DeleteUser", err)
	}

	if err := authService.CheckPassword(user.Password, req.GetPassword()); err != nil {
		return nil, helper.RespondWithErrorGRPC(ctx, codes.InvalidArgument, "password is not correct: DeleteUser", err)
	}

	if req.GetVerifyMessage() != "SUBMIT" {
		return nil, helper.RespondWithErrorGRPC(ctx, codes.InvalidArgument, "you must submit 'SUBMIT' to delete your account", nil)
	}

	if err = s.db.DeleteUser(ctx, userID); err != nil {
		return nil, helper.RespondWithErrorGRPC(ctx, codes.Internal, "can't delete user from db: DeleteUser", err)
	}

	return &pb.DeleteUserResponse{
		Status: "success",
	}, nil
}

func (s *server) SendVerificationCode(
	ctx context.Context,
	req *pb.SendVerificationCodeRequest,
) (*pb.SendVerificationCodeResponse, error) {

	accessToken, err := postService.GetBearerTokenFromGrpc(ctx)
	if err != nil {
		return nil, helper.RespondWithErrorGRPC(ctx, codes.InvalidArgument, "can't get authorization token from header: SendVerificationCode", err)
	}

	userID, err := postService.ValidateJWT(accessToken, s.tokenSecret)
	if err != nil {
		return nil, helper.RespondWithErrorGRPC(ctx, codes.Unauthenticated, "you are not allowed to reset password: SendVerificationCode", err)
	}

	user, err := s.db.GetUserById(ctx, userID)
	if err != nil {
		return nil, helper.RespondWithErrorGRPC(ctx, codes.Internal, "can't get user from db: SendVerificationCode", err)
	}

	verificationCode, err := authService.GenerateVerificationCode()
	if err != nil {
		return nil, helper.RespondWithErrorGRPC(ctx, codes.Internal, "can't generate verification code: SendVerificationCode", err)
	}

	sendResetVerificationCodeParams := database.SendResetVerificationCodeParams{
		ID:               userID,
		VerificationCode: verificationCode,
	}

	err = s.db.SendResetVerificationCode(ctx, sendResetVerificationCodeParams)
	if err != nil {
		return nil, helper.RespondWithErrorGRPC(ctx, codes.Internal, "can't send verification code: SendVerificationCode", err)
	}

	err = authService.SendVerificationEmail(user.Email, s.email, s.emailSecret, verificationCode)
	if err != nil {
		return nil, helper.RespondWithErrorGRPC(ctx, codes.Internal, "can't send verification email: SendVerificationCode", err)
	}

	return &pb.SendVerificationCodeResponse{
		Status: "Verification code sent",
	}, nil
}

func (s *server) ResetPassword(ctx context.Context, req *pb.ResetPasswordRequest) (*pb.ResetPasswordResponse, error) {
	accessToken, err := postService.GetBearerTokenFromGrpc(ctx)
	if err != nil {
		return nil, helper.RespondWithErrorGRPC(ctx, codes.InvalidArgument, "can't get authorization token from header: ResetPassword", err)
	}

	userID, err := postService.ValidateJWT(accessToken, s.tokenSecret)
	if err != nil {
		return nil, helper.RespondWithErrorGRPC(ctx, codes.Unauthenticated, "you are not allowed to reset password: ResetPassword", err)
	}

	user, err := s.db.GetUserById(ctx, userID)
	if err != nil {
		return nil, helper.RespondWithErrorGRPC(ctx, codes.Internal, "can't get user from db: ResetPassword", err)
	}

	if req.GetVerificationCode() != user.VerificationCode {
		return nil, helper.RespondWithErrorGRPC(ctx, codes.InvalidArgument, "verification code is not correct", nil)
	}

	err = s.db.VerifyVerificationCode(ctx, userID)
	if err != nil {
		return nil, helper.RespondWithErrorGRPC(ctx, codes.Internal, "can't verify verification code: ResetPassword", err)
	}

	newPassword, err := authService.HashPassword(req.NewPassword)
	if err != nil {
		return nil, helper.RespondWithErrorGRPC(ctx, codes.Internal, "can't hash new password: ResetPassword", err)
	}

	resetPasswordParams := database.ResetPasswordParams{
		ID:       userID,
		Password: newPassword,
	}

	err = s.db.ResetPassword(ctx, resetPasswordParams)
	if err != nil {
		return nil, helper.RespondWithErrorGRPC(ctx, codes.Internal, "can't reset password: ResetPassword", err)
	}

	return &pb.ResetPasswordResponse{
		Status: "Password changed successfully",
	}, nil
}

func (s *server) DeleteAllUsers(ctx context.Context, req *pb.DeleteAllUsersRequest) (*pb.DeleteAllUsersResponse, error) {
	err := s.db.DeleteAllUsers(ctx)
	if err != nil {
		return nil, helper.RespondWithErrorGRPC(ctx, codes.Internal, "can't delete user from db: DeleteAllUsers", err)
	}

	return &pb.DeleteAllUsersResponse{
		Status: "users deleted successfully",
	}, nil
}
