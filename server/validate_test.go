package server

import (
	"testing"

	"github.com/lileio/auth_service"
	"github.com/stretchr/testify/assert"
	context "golang.org/x/net/context"
)

func TestValidate(t *testing.T) {
	token, err := NewToken("someaccountid")
	assert.Nil(t, err)

	ctx := context.Background()
	req := &auth_service.ValidateRequest{Token: token}

	_, err = cli.Validate(ctx, req)
	assert.Nil(t, err)
}
