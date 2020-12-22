package commons

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	//port = flag.Int("port", 50051, "the port to serve on")

	ErrMissingMetadata = status.Errorf(codes.InvalidArgument, "missing metadata")
	ErrInvalidToken    = status.Errorf(codes.Unauthenticated, "invalid token")
	ErrExpiredToken    = status.Errorf(codes.Unauthenticated, "expired token")
	UnimplementedError = status.Errorf(codes.Unimplemented, "Not done yet")
	UserNotFound       = status.Errorf(codes.NotFound, "User not found")
	AlreadyExist       = status.Errorf(codes.AlreadyExists, "Already exist")
	UserAlreadyExist   = status.Errorf(codes.AlreadyExists, "User already exist")
	NotAuthorized      = status.Errorf(codes.PermissionDenied, "Not Authorrized")

	EmailAlreadyUsed  = status.Errorf(codes.AlreadyExists, "Email already used")
	MobileAlreadyUsed = status.Errorf(codes.AlreadyExists, "Mobile already used")

	OTPWrong = status.Errorf(codes.FailedPrecondition, "Wrong OTP")

	NotFound = status.Errorf(codes.NotFound, "Not found")
)
