package api

import (
	pb "github.com/Mujib-Ahasan/Rampaz/proto"
	"google.golang.org/grpc"
)

func (s *K8SServer) GetNodeRealTimeStats(req *pb.NodeRequest, stream grpc.ServerStreamingServer[pb.NodeStatsResponse]) error {

	return s.NodeMetService.StreamNodeStats(
		stream.Context(),
		req.NodeName,
		stream.Send,
	)
}
