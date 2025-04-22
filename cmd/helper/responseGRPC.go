package helper

import (
	"context"
	"encoding/json"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RespondWithErrorGRPC creates a standardized gRPC error response with the given code and message.
// It also logs the error if provided.
func RespondWithErrorGRPC(ctx context.Context, code codes.Code, msg string, err error) error {
	if err != nil {
		log.Println(err)
	}

	if code > codes.Internal { // 5XX equivalent in gRPC
		log.Printf("Responding with 5XX gRPC error: %s", msg)
	}

	type errorResponse struct {
		UserServiceError string `json:"error"`
	}

	jsonBytes, err := json.Marshal(errorResponse{UserServiceError: msg})
	if err != nil {
		log.Printf("Error marshalling error JSON: %s", err)
		return status.Errorf(codes.Internal, "Failed to marshal error response")
	}

	log.Printf("UserServiceError: %s, Code: %s", string(jsonBytes), code.String()) // Log the error
	return status.Error(code, msg)
}
