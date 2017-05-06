package server

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/lileio/account_service"
	"github.com/lileio/auth_service"
	context "golang.org/x/net/context"
)

func (s Server) Authenticate(ctx context.Context, r *auth_service.AuthRequest) (*auth_service.AuthResponse, error) {
	as := AccountService()
	req := account_service.AuthenticateByEmailRequest{
		Email:    r.Email,
		Password: r.Password,
	}
	res, err := as.AuthenticateByEmail(ctx, &req)
	if err != nil {
		return nil, grpc.Errorf(
			codes.PermissionDenied,
			"authentication incorrect: %s",
			err)
	}

	t, err := NewToken(res.Id)
	if err != nil {
		return nil, grpc.Errorf(
			codes.Internal,
			"jwt signing error: %s",
			err)
	}

	ar := auth_service.AuthResponse{
		Token:   t,
		Account: res,
	}

	return &ar, nil
}
