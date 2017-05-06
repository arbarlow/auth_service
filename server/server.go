package server

import (
	"errors"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/lileio/account_service"
	"github.com/lileio/auth_service"
	"github.com/lileio/lile"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type Server struct {
	auth_service.AuthServiceServer
}

type JWTClaims struct {
	AccountID string `json:"account_id"`
	jwt.StandardClaims
}

var (
	as         account_service.AccountServiceClient
	signingKey string
)

func init() {
	key := os.Getenv("SIGNING_KEY")
	if key == "" {
		logrus.Fatal("no JWT signing key provided")
	}

	signingKey = key
}

func NewServer() *lile.Server {
	s := &Server{}

	impl := func(g *grpc.Server) {
		auth_service.RegisterAuthServiceServer(g, s)
	}

	return lile.NewServer(
		lile.Name("auth_service"),
		lile.Implementation(impl),
		lile.Publishers(map[string]string{}),
	)
}

func NewToken(account_id string) (string, error) {
	if account_id == "" {
		return "", errors.New("no account_id provided")
	}

	claims := JWTClaims{
		account_id,
		jwt.StandardClaims{
			ExpiresAt: 15000,
			Issuer:    "auth_token_service",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(signingKey))
}

func AccountService() account_service.AccountServiceClient {
	if as != nil {
		return as
	}

	addr := os.Getenv("ACCOUNT_SERVICE_ADDR")
	if addr == "" {
		addr = "account_service"
	}

	t := opentracing.GlobalTracer()

	conn, err := grpc.Dial(
		addr,
		grpc.WithInsecure(),
		grpc.WithTimeout(1*time.Second),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(t)),
	)
	if err != nil {
		logrus.Warnf("account service connection error: %s", err)
	}

	as = account_service.NewAccountServiceClient(conn)
	return as

}
