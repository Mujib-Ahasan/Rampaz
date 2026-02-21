package api

import (
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
		return err
	}

	s := grpc.NewServer()
	pb.RegisterK8SInfoServer(s, svc)
	reflection.Register(s)

	return s.Serve(lis)
}
