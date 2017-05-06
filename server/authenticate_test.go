package server

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lileio/account_service"
	"github.com/lileio/auth_service"
	"github.com/lileio/auth_service/server/mock_account_service"
	"github.com/stretchr/testify/assert"
	context "golang.org/x/net/context"
)

func TestAuthenticate(t *testing.T) {
	email := "somecorrect@email"
	acc := account_service.Account{Id: "GUID-0001", Email: email}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock_account_service.NewMockAccountServiceClient(ctrl)
	m.EXPECT().AuthenticateByEmail(gomock.Any(), gomock.Any()).
		Times(1).Return(&acc, nil)
	as = m

	ctx := context.Background()
	req := &auth_service.AuthRequest{
		Email:    email,
		Password: "correctpassword",
	}

	res, err := cli.Authenticate(ctx, req)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.NotEmpty(t, res.Token)
}

func TestAuthenticateFailure(t *testing.T) {
	email := "somecorrect@email"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock_account_service.NewMockAccountServiceClient(ctrl)
	m.EXPECT().AuthenticateByEmail(gomock.Any(), gomock.Any()).
		Times(1).Return(nil, errors.New("something failed"))
	as = m

	ctx := context.Background()
	req := &auth_service.AuthRequest{
		Email:    email,
		Password: "correctpassword",
	}

	res, err := cli.Authenticate(ctx, req)
	assert.NotNil(t, err)
	assert.Nil(t, res)
}
