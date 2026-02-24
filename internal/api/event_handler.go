package api

import (
	"github.com/Mujib-Ahasan/Rampaz/internal/metrics"
	pb "github.com/Mujib-Ahasan/Rampaz/proto"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *K8SServer) StreamEvents(_ *emptypb.Empty, stream grpc.ServerStreamingServer[pb.EventResponse]) error {
	endpoint := "stream_events"
	status := "success"

	timer := prometheus.NewTimer(
		metrics.RequestLatency.WithLabelValues(endpoint),
	)

	metrics.ActiveStreams.
		WithLabelValues(endpoint).
		Inc()
	defer metrics.ActiveStreams.
		WithLabelValues(endpoint).
		Dec()

	defer func() {
		timer.ObserveDuration()
		metrics.APIRequests.
			WithLabelValues(endpoint, status).
			Inc()
	}()

	ctx := stream.Context()

	events, err := s.EventService.StreamEvents(ctx)
	if err != nil {
		status = "error"
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
				status = "error"
				return err
			}

			metrics.StreamMessagesSent.
				WithLabelValues(endpoint).
				Inc()
		}
	}
}
