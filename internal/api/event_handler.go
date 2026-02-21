package api

import (
	pb "github.com/Mujib-Ahasan/Rampaz/proto"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *K8SServer) StreamEvents(_ *emptypb.Empty, stream grpc.ServerStreamingServer[pb.EventResponse]) error {

	ctx := stream.Context()

	events, err := s.EventService.StreamEvents(ctx)
	if err != nil {
		return err
	}

	for {
		select {

		case <-ctx.Done():
			return nil

		case ev, ok := <-events:
			if !ok {
				return nil
			}

			resp := &pb.EventResponse{
				Type:           ev.Type,
				Reason:         ev.Reason,
				Message:        ev.Message,
				InvolvedObject: ev.InvolvedObject.Kind + "/" + ev.InvolvedObject.Name,
			}

			if err := stream.Send(resp); err != nil {
				return err
			}
		}
	}
}
