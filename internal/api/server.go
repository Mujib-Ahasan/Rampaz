package api

import (
	"fmt"
	"net"

	pb "github.com/Mujib-Ahasan/Rampaz/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type ServerDeps struct {
	Kube    any
	Metrics any
}

func StartGRPC(addr string, svc pb.K8SInfoServer) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("start grpc server: failed to listen on %q: %w", addr, err)
	}

	s := grpc.NewServer()
	pb.RegisterK8SInfoServer(s, svc)
	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		return fmt.Errorf("start grpc server: serve failed on %q: %w", addr, err)
	}

	return nil
}
