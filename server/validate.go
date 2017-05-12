package server

import (
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/lileio/auth_service"
	context "golang.org/x/net/context"
)

func (s Server) Validate(ctx context.Context, r *auth_service.ValidateRequest) (*empty.Empty, error) {
	err := ValidateToken(r.Token)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}
