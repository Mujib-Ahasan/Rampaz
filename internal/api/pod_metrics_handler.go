package api

import (
	pb "github.com/Mujib-Ahasan/Rampaz/proto"
	"google.golang.org/grpc"
)

func (s *K8SServer) GetPodStats(req *pb.PodRequest, stream grpc.ServerStreamingServer[pb.PodStatsResponse]) error {

	return s.PodMetService.StreamPodStats(
		stream.Context(),
		req.Namespace,
		stream.Send,
	)
}
