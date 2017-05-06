package server

import (
	"os"
	"testing"

	"google.golang.org/grpc"

	"github.com/lileio/auth_service"
	"github.com/lileio/lile"
)

var s = Server{}
var cli auth_service.AuthServiceClient

func TestMain(m *testing.M) {
	impl := func(g *grpc.Server) {
		auth_service.RegisterAuthServiceServer(g, s)
	}

	addr, serve := lile.NewTestServer(
		lile.Name("auth_service_test"),
		lile.Implementation(impl),
	)

	go serve()

	cli = auth_service.NewAuthServiceClient(lile.TestConn(addr))

	os.Exit(m.Run())
}
