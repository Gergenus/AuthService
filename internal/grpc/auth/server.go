package auth

import (
	"context"
	"errors"

	"github.com/Gergenus/AuthService/internal/repository"
	authv1 "github.com/Gergenus/Protobuf/gen/go/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	SignIn(ctx context.Context, email string, password string) (token string, err error)
	GetUser(ctx context.Context, username string) (userID int64, err error)
	RegisterNewUser(ctx context.Context, username string, email string, password string) (userID int64, err error)
}

type serverAPI struct {
	authv1.UnimplementedAuthServer
	auth Auth
}

func Register(gRPC *grpc.Server, auth Auth) {
	authv1.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
}

func (s *serverAPI) SignIn(ctx context.Context, req *authv1.SignInRequest) (*authv1.SignInResponse, error) {
	if req.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	if req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	token, err := s.auth.SignIn(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &authv1.SignInResponse{
		Token: token,
	}, nil
}

func (s *serverAPI) GetUserId(ctx context.Context, req *authv1.GetUserIdRequest) (*authv1.GetUserIdResponse, error) {
	if req.GetUsername() == "" {
		return nil, status.Error(codes.InvalidArgument, "username is required")
	}

	id, err := s.auth.GetUser(ctx, req.GetUsername())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &authv1.GetUserIdResponse{
		UserId: id,
	}, nil
}

func (s *serverAPI) SignUp(ctx context.Context, req *authv1.SignUpRequest) (*authv1.SignUpResponse, error) {
	if req.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	if req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	id, err := s.auth.RegisterNewUser(ctx, req.GetUsername(), req.GetEmail(), req.GetPassword())
	if err != nil {
		if errors.Is(err, repository.ErrUserEXists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &authv1.SignUpResponse{
		UserId: int64(id),
	}, nil
}

func (s *serverAPI) mustEmbedUnimplementedAuthServer() {

}
