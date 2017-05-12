package server

import (
	"errors"
	"fmt"
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
	as          account_service.AccountServiceClient
	signingKey  string
	tokenExpiry time.Duration
)

func init() {
	signingKey := os.Getenv("SIGNING_TOKEN")
	if signingKey == "" {
		logrus.Fatal("no JWT signing token provided (SIGNING_TOKEN)")
	}

	envDuration := os.Getenv("TOKEN_EXPIRY")
	if envDuration == "" {
		envDuration = "48h"
	}

	d, err := time.ParseDuration(envDuration)
	if err != nil {
		logrus.Fatal(err)
	}
	tokenExpiry = d
}

func NewServer() *lile.Server {
	s := &Server{}

	impl := func(g *grpc.Server) {
		auth_service.RegisterAuthServiceServer(g, s)
	}

	return lile.NewServer(
		lile.Name("auth_service"),
		lile.Implementation(impl),
	)
}

func NewToken(account_id string) (string, error) {
	if account_id == "" {
		return "", errors.New("no account_id provided")
	}

	claims := JWTClaims{
		account_id,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenExpiry).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString([]byte(signingKey))
}

func ValidateToken(token string) error {
	_, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(signingKey), nil
	})

	return err
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
