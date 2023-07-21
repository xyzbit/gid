package server

import (
	"log"
	"net"

	"github.com/google/wire"
	v1 "github.com/xyzbit/gid/api/v1"
	"github.com/xyzbit/gid/core/conf"
	"google.golang.org/grpc"
)

var ProviderSet = wire.NewSet(NewGeneratorSvc, NewGrpcServer)

type GrpcServer struct {
	Addr string
	srv  *grpc.Server
	gs   *GeneratorSvc
}

// NewServer new gRPC server instance
func NewGrpcServer(gs *GeneratorSvc, c *conf.Server) *GrpcServer {
	if c == nil && c.Grpc == nil {
		panic("server.grpc config is nil")
	}
	return &GrpcServer{
		Addr: c.Grpc.Addr,
		gs:   gs,
	}
}

// Start start gRPC server
func (s *GrpcServer) Start() error {
	listenAddr := s.Addr

	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}

	srv := grpc.NewServer()
	s.srv = srv
	v1.RegisterGeneratorServer(srv, s.gs)
	v1.RegisterMannagerServer(srv, s.gs)

	log.Printf("gRPC server listening on %s", listenAddr)
	if err := srv.Serve(lis); err != nil {
		return err
	}
	return nil
}

// Stop stop gRPC server
func (s *GrpcServer) Stop() {
	s.srv.GracefulStop()
}
