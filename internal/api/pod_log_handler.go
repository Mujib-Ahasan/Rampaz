package api

import (
	"context"
	"encoding/json"
	"strings"

	pb "github.com/Mujib-Ahasan/Rampaz/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var allowedLevels = map[string]bool{
	"error":   true,
	"warn":    true,
	"warning": true,
	"fatal":   true,
	"panic":   true,
}

const maxEntries = 50

type jsonLog struct {
	Level string `json:"level"`
	Ts    string `json:"ts"`
	Msg   string `json:"msg"`
	Error string `json:"error"`
}

func (s *K8SServer) GetPodLogs(ctx context.Context, req *pb.PodLogsRequest) (*pb.PodLogsResponse, error) {
	if req.Namespace == "" {
		return nil, status.Error(codes.InvalidArgument, "namespace is required")
	}
	if req.PodName == "" {
		return nil, status.Error(codes.InvalidArgument, "pod name is required")
	}

	logs, err := s.PodLogService.GetPodLogs(ctx, req.Namespace, req.PodName, req.TailLines)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get pod logs: %v", err)
	}

	return &pb.PodLogsResponse{
		Namespace: req.Namespace,
		PodName:   req.PodName,
		Entries:   FilterAndStructureLogs(logs),
	}, nil
}

func FilterAndStructureLogs(raw string) []*pb.PodLogEntry {
	lines := strings.Split(raw, "\n")
	entries := make([]*pb.PodLogEntry, 0, maxEntries)

	for i := len(lines) - 1; i >= 0; i-- {
		if len(entries) >= maxEntries {
			break
		}

		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}

		// Try parsing JSON logs
		var jl jsonLog
		if err := json.Unmarshal([]byte(line), &jl); err == nil {
			level := strings.ToLower(jl.Level)

			if !allowedLevels[level] {
				continue
			}

			entries = append(entries, &pb.PodLogEntry{
				Level:     level,
				Timestamp: jl.Ts,
				Message:   jl.Msg,
				Error:     jl.Error,
			})
			continue
		}

		// Fallback for non-JSON logs
		lower := strings.ToLower(line)
		if strings.Contains(lower, "error") ||
			strings.Contains(lower, "warn") ||
			strings.Contains(lower, "fatal") ||
			strings.Contains(lower, "panic") {

			entries = append(entries, &pb.PodLogEntry{
				Level:   "error", // fallback assumption
				Message: line,
			})
		}
	}

	// Reverse to maintain chronological order
	for i, j := 0, len(entries)-1; i < j; i, j = i+1, j-1 {
		entries[i], entries[j] = entries[j], entries[i]
	}

	return entries
}
