package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	pb "github.com/Mujib-Ahasan/Rampaz/proto"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	defaultToolTimeout = 10 * time.Second
	maxStreamItems     = 20
	podLogTailLines    = 50
)

type ToolExecutor struct {
	client pb.K8SInfoClient
}

func NewToolExecutor(client pb.K8SInfoClient) *ToolExecutor {
	return &ToolExecutor{
		client: client,
	}
}

func (e *ToolExecutor) Execute(ctx context.Context, call ToolCall) (ToolResult, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultToolTimeout)
	defer cancel()

	switch call.Name {
	case "list_pods":
		return e.ListPods(ctx, call.Args)

	case "get_node_stats":
		return e.GetNodeStats(ctx, call.Args)

	case "get_pod_stats":
		return e.GetPodStats(ctx, call.Args)

	case "get_recent_events":
		return e.GetRecentEvents(ctx, call.Args)

	case "list_deployments":
		return e.ListDeployments(ctx, call.Args)

	case "list_replicasets":
		return e.ListReplicaSets(ctx, call.Args)

	case "list_statefulsets":
		return e.ListStatefulSets(ctx, call.Args)

	case "list_daemonsets":
		return e.ListDaemonSets(ctx, call.Args)

	case "list_jobs":
		return e.ListJobs(ctx, call.Args)

	case "list_cronjobs":
		return e.ListCronJobs(ctx, call.Args)

	case "list_services":
		return e.ListServices(ctx, call.Args)

	case "list_namespaces":
		return e.ListNamespaces(ctx, call.Args)

	case "list_pvcs":
		return e.ListPVCs(ctx, call.Args)

	case "list_pvs":
		return e.ListPVs(ctx, call.Args)

	case "list_nodes":
		return e.ListNodes(ctx, call.Args)

	case "list_network_policies":
		return e.ListNetworkPolicies(ctx, call.Args)

	case "get_namespace_summary":
		return e.GetNamespaceSummary(ctx, call.Args)

	case "get_cluster_overview":
		return e.GetClusterOverview(ctx, call.Args)

	case "get_workloads_by_health":
		return e.GetWorkloadsByHealth(ctx, call.Args)

	case "get_namespace_metrics":
		return e.GetNamespaceMetrics(ctx, call.Args)

	case "get_node_resource_allocation":
		return e.GetNodeResourceAllocation(ctx, call.Args)
	case "get_pod_logs":
		return e.GetPodLogs(ctx, call.Args)

	default:
		return ToolResult{}, fmt.Errorf("unknown tool: %s", call.Name)
	}
}

func (e *ToolExecutor) ListPods(ctx context.Context, args map[string]any) (ToolResult, error) {
	namespace := getStringArg(args, "namespace")
	labelSelector := getStringArg(args, "label_selector")

	resp, err := e.client.ListPods(ctx, &pb.NamespaceRequest{
		Namespace:     namespace,
		LabelSelector: labelSelector,
	})
	return result("list_pods", resp, err)
}

func (e *ToolExecutor) GetNodeStats(ctx context.Context, args map[string]any) (ToolResult, error) {
	nodeName := getStringArg(args, "nodeName")
	if nodeName == "" {
		return ToolResult{}, fmt.Errorf("nodeName is required")
	}

	resp, err := e.client.GetNodeStats(ctx, &pb.NodeRequest{
		NodeName: nodeName,
	})
	return result("get_node_stats", resp, err)
}

func (e *ToolExecutor) GetPodStats(ctx context.Context, args map[string]any) (ToolResult, error) {
	namespace := getStringArg(args, "namespace")
	limit := getIntArg(args, "limit", maxStreamItems)

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	stream, err := e.client.GetPodStats(ctx, &pb.PodRequest{
		Namespace: namespace,
	})
	if err != nil {
		return ToolResult{}, err
	}

	items := make([]*pb.PodStatsResponse, 0, limit)

	for len(items) < limit {
		item, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			if ctx.Err() != nil {
				break
			}
			return ToolResult{}, err
		}

		items = append(items, item)
	}

	return ToolResult{
		Name: "get_pod_stats",
		Data: items,
	}, nil
}

func (e *ToolExecutor) GetRecentEvents(ctx context.Context, args map[string]any) (ToolResult, error) {
	namespace := getStringArg(args, "namespace")
	limit := getIntArg(args, "limit", maxStreamItems)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	stream, err := e.client.StreamEvents(ctx, &pb.NamespaceRequest{
		Namespace: namespace,
	})

	if err != nil {
		return ToolResult{}, err
	}

	events := make([]*pb.EventResponse, 0, limit)

	for len(events) < limit {
		event, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			if ctx.Err() != nil {
				break
			}
			return ToolResult{}, err
		}

		events = append(events, event)
	}

	return ToolResult{
		Name: "get_recent_events",
		Data: events,
	}, nil
}

func (e *ToolExecutor) GetPodLogs(ctx context.Context, args map[string]any) (ToolResult, error) {
	namespace := getStringArg(args, "namespace")
	podName := getStringArg(args, "pod_name")

	resp, err := e.client.GetPodLogs(ctx, &pb.PodLogsRequest{
		Namespace: namespace,
		PodName:   podName,
		TailLines: podLogTailLines,
	})
	return result("get_pod_logs", resp, err)
}

func (e *ToolExecutor) ListDeployments(ctx context.Context, args map[string]any) (ToolResult, error) {
	return e.listWorkload(ctx, "list_deployments", args, e.client.ListDeployments)
}

func (e *ToolExecutor) ListReplicaSets(ctx context.Context, args map[string]any) (ToolResult, error) {
	return e.listWorkload(ctx, "list_replicasets", args, e.client.ListReplicaSets)
}

func (e *ToolExecutor) ListStatefulSets(ctx context.Context, args map[string]any) (ToolResult, error) {
	return e.listWorkload(ctx, "list_statefulsets", args, e.client.ListStatefulSets)
}

func (e *ToolExecutor) ListDaemonSets(ctx context.Context, args map[string]any) (ToolResult, error) {
	return e.listWorkload(ctx, "list_daemonsets", args, e.client.ListDaemonSets)
}

func (e *ToolExecutor) ListJobs(ctx context.Context, args map[string]any) (ToolResult, error) {
	return e.listWorkload(ctx, "list_jobs", args, e.client.ListJobs)
}

func (e *ToolExecutor) ListCronJobs(ctx context.Context, args map[string]any) (ToolResult, error) {
	return e.listWorkload(ctx, "list_cronjobs", args, e.client.ListCronJobs)
}

func (e *ToolExecutor) listWorkload(
	ctx context.Context,
	name string,
	args map[string]any,
	fn func(context.Context, *pb.NamespaceRequest, ...grpc.CallOption) (*pb.WorkloadListResponse, error),
) (ToolResult, error) {
	namespace := getStringArg(args, "namespace")
	labelSelector := getStringArg(args, "label_selector")

	resp, err := fn(ctx, &pb.NamespaceRequest{
		Namespace:     namespace,
		LabelSelector: labelSelector,
	})

	return result(name, resp, err)
}

func (e *ToolExecutor) ListServices(ctx context.Context, args map[string]any) (ToolResult, error) {
	namespace := getStringArg(args, "namespace")
	labelSelector := getStringArg(args, "label_selector")

	resp, err := e.client.ListServices(ctx, &pb.NamespaceRequest{
		Namespace:     namespace,
		LabelSelector: labelSelector,
	})
	return result("list_services", resp, err)
}

func (e *ToolExecutor) ListNamespaces(ctx context.Context, args map[string]any) (ToolResult, error) {
	resp, err := e.client.ListNamespaces(ctx, &emptypb.Empty{})
	return result("list_namespaces", resp, err)
}

func (e *ToolExecutor) ListPVCs(ctx context.Context, args map[string]any) (ToolResult, error) {
	namespace := getStringArg(args, "namespace")
	labelSelector := getStringArg(args, "label_selector")

	resp, err := e.client.ListPVCs(ctx, &pb.NamespaceRequest{
		Namespace:     namespace,
		LabelSelector: labelSelector,
	})
	return result("list_pvcs", resp, err)
}

func (e *ToolExecutor) ListPVs(ctx context.Context, args map[string]any) (ToolResult, error) {
	resp, err := e.client.ListPVs(ctx, &emptypb.Empty{})
	return result("list_pvs", resp, err)
}

func (e *ToolExecutor) ListNodes(ctx context.Context, args map[string]any) (ToolResult, error) {
	resp, err := e.client.ListNodes(ctx, &emptypb.Empty{})
	return result("list_nodes", resp, err)
}

func (e *ToolExecutor) ListNetworkPolicies(ctx context.Context, args map[string]any) (ToolResult, error) {
	namespace := getStringArg(args, "namespace")
	labelSelector := getStringArg(args, "label_selector")

	resp, err := e.client.ListNetworkPolicies(ctx, &pb.NamespaceRequest{
		Namespace:     namespace,
		LabelSelector: labelSelector,
	})
	return result("list_network_policies", resp, err)
}

func (e *ToolExecutor) GetNamespaceSummary(ctx context.Context, args map[string]any) (ToolResult, error) {
	namespace := getStringArg(args, "namespace")
	if namespace == "" {
		return ToolResult{}, fmt.Errorf("namespace is required")
	}

	resp, err := e.client.GetNamespaceSummary(ctx, &pb.NamespaceRequest{
		Namespace: namespace,
	})
	return result("get_namespace_summary", resp, err)
}

func (e *ToolExecutor) GetClusterOverview(ctx context.Context, args map[string]any) (ToolResult, error) {
	resp, err := e.client.GetClusterOverview(ctx, &emptypb.Empty{})
	return result("get_cluster_overview", resp, err)
}

func (e *ToolExecutor) GetWorkloadsByHealth(ctx context.Context, args map[string]any) (ToolResult, error) {
	health := strings.ToUpper(getStringArg(args, "health"))
	if health == "" {
		return ToolResult{}, fmt.Errorf("health is required")
	}

	namespace := getStringArg(args, "namespace")

	resp, err := e.client.GetWorkloadsByHealth(ctx, &pb.WorkloadHealthRequest{
		Namespace: namespace,
		Health:    health,
	})
	return result("get_workloads_by_health", resp, err)
}

func (e *ToolExecutor) GetNamespaceMetrics(ctx context.Context, args map[string]any) (ToolResult, error) {
	namespace := getStringArg(args, "namespace")
	if namespace == "" {
		return ToolResult{}, fmt.Errorf("namespace is required")
	}

	resp, err := e.client.GetNamespaceMetrics(ctx, &pb.NamespaceMetricsRequest{
		Namespace: namespace,
	})
	return result("get_namespace_metrics", resp, err)
}

func (e *ToolExecutor) GetNodeResourceAllocation(ctx context.Context, args map[string]any) (ToolResult, error) {
	nodeName := getStringArg(args, "nodeName")
	if nodeName == "" {
		return ToolResult{}, fmt.Errorf("nodeName is required")
	}

	resp, err := e.client.GetNodeResourceAllocation(ctx, &pb.NodeRequest{
		NodeName: nodeName,
	})
	return result("get_node_resource_allocation", resp, err)
}

func result(name string, data any, err error) (ToolResult, error) {
	if err != nil {
		return ToolResult{}, err
	}

	return ToolResult{
		Name: name,
		Data: data,
	}, nil
}

func getStringArg(args map[string]any, key string) string {
	if args == nil {
		return ""
	}

	v, ok := args[key]
	if !ok || v == nil {
		return ""
	}

	switch val := v.(type) {
	case string:
		return strings.TrimSpace(val)
	default:
		return strings.TrimSpace(fmt.Sprint(val))
	}
}

func getIntArg(args map[string]any, key string, fallback int) int {
	if args == nil {
		return fallback
	}

	v, ok := args[key]
	if !ok || v == nil {
		return fallback
	}

	switch val := v.(type) {
	case int:
		return val
	case int32:
		return int(val)
	case int64:
		return int(val)
	case float64:
		return int(val)
	case json.Number:
		i, err := val.Int64()
		if err == nil {
			return int(i)
		}
	}

	return fallback
}
