package ai

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	pb "github.com/Mujib-Ahasan/Rampaz/proto"
)

const (
	maxContextChars       = 12000
	maxPodsInContext      = 30
	maxEventsInContext    = 30
	maxWorkloadsInContext = 30
	maxPodStatsInContext  = 30
)

func BuildClusterContext(results []ToolResult) string {
	var b strings.Builder

	b.WriteString("Cluster context generated from Rampaz tool results.\n")
	b.WriteString("Use only the facts below. Do not invent cluster data.\n\n")

	for _, result := range results {
		if result.Error != "" {
			writeToolError(&b, result)
			continue
		}

		switch result.Name {
		case "get_cluster_overview":
			writeClusterOverview(&b, result.Data)

		case "get_namespace_summary":
			writeNamespaceSummary(&b, result.Data)

		case "get_namespace_metrics":
			writeNamespaceMetrics(&b, result.Data)

		case "get_workloads_by_health":
			writeWorkloadsByHealth(&b, result.Data)

		case "list_pods":
			writePods(&b, result.Data)

		case "get_pod_stats":
			writePodStats(&b, result.Data)

		case "get_recent_events":
			writeRecentEvents(&b, result.Data)

		case "list_nodes":
			writeNodes(&b, result.Data)

		case "get_node_stats":
			writeNodeStats(&b, result.Data)

		case "get_node_resource_allocation":
			writeNodeResourceAllocation(&b, result.Data)

		default:
			writeFallback(&b, result)
		}

		b.WriteString("\n")
	}

	out := b.String()
	if len(out) > maxContextChars {
		out = out[:maxContextChars] + "\n\n[Context truncated]\n"
	}

	return out
}

func writeToolError(b *strings.Builder, result ToolResult) {
	b.WriteString("Tool Error:\n")
	b.WriteString(fmt.Sprintf("- Tool: %s\n", result.Name))
	b.WriteString(fmt.Sprintf("- Error: %s\n\n", result.Error))
}

func writeClusterOverview(b *strings.Builder, data any) {
	resp, ok := data.(*pb.ClusterOverviewResponse)
	if !ok {
		writeRawData(b, "Cluster Overview", data)
		return
	}

	b.WriteString("Cluster Overview:\n")
	b.WriteString(fmt.Sprintf("- Nodes: %d\n", resp.GetNodes()))
	b.WriteString(fmt.Sprintf("- Namespaces: %d\n", resp.GetNamespaces()))
	b.WriteString(fmt.Sprintf("- Pods: %d\n", resp.GetPods()))
	b.WriteString(fmt.Sprintf("- Deployments: %d\n", resp.GetDeployments()))
	b.WriteString(fmt.Sprintf("- ReplicaSets: %d\n", resp.GetReplicasets()))
	b.WriteString(fmt.Sprintf("- StatefulSets: %d\n", resp.GetStatefulsets()))
	b.WriteString(fmt.Sprintf("- DaemonSets: %d\n", resp.GetDaemonsets()))
	b.WriteString(fmt.Sprintf("- Jobs: %d\n", resp.GetJobs()))
	b.WriteString(fmt.Sprintf("- CronJobs: %d\n", resp.GetCronjobs()))
	b.WriteString(fmt.Sprintf("- Services: %d\n", resp.GetServices()))
	b.WriteString(fmt.Sprintf("- PVCs: %d\n", resp.GetPersistentVolumeClaims()))
	b.WriteString(fmt.Sprintf("- NetworkPolicies: %d\n", resp.GetNetworkPolicies()))
}

func writeNamespaceSummary(b *strings.Builder, data any) {
	resp, ok := data.(*pb.NamespaceSummaryResponse)
	if !ok {
		writeRawData(b, "Namespace Summary", data)
		return
	}

	b.WriteString("Namespace Summary:\n")
	b.WriteString(fmt.Sprintf("- Namespace: %s\n", resp.GetNamespace()))
	b.WriteString(fmt.Sprintf("- Pods: %d\n", resp.GetPods()))
	b.WriteString(fmt.Sprintf("- Deployments: %d\n", resp.GetDeployments()))
	b.WriteString(fmt.Sprintf("- ReplicaSets: %d\n", resp.GetReplicasets()))
	b.WriteString(fmt.Sprintf("- StatefulSets: %d\n", resp.GetStatefulsets()))
	b.WriteString(fmt.Sprintf("- DaemonSets: %d\n", resp.GetDaemonsets()))
	b.WriteString(fmt.Sprintf("- Jobs: %d\n", resp.GetJobs()))
	b.WriteString(fmt.Sprintf("- CronJobs: %d\n", resp.GetCronjobs()))
	b.WriteString(fmt.Sprintf("- Services: %d\n", resp.GetServices()))
	b.WriteString(fmt.Sprintf("- PVCs: %d\n", resp.GetPersistentVolumeClaims()))
	b.WriteString(fmt.Sprintf("- NetworkPolicies: %d\n", resp.GetNetworkPolicies()))
}

func writeNamespaceMetrics(b *strings.Builder, data any) {
	resp, ok := data.(*pb.NamespaceMetricsResponse)
	if !ok {
		writeRawData(b, "Namespace Metrics", data)
		return
	}

	b.WriteString("Namespace Metrics:\n")
	b.WriteString(fmt.Sprintf("- Namespace: %s\n", resp.GetNamespace()))
	b.WriteString(fmt.Sprintf("- Total pods: %d\n", resp.GetTotalPods()))
	b.WriteString(fmt.Sprintf("- Running pods: %d\n", resp.GetRunningPods()))
	b.WriteString(fmt.Sprintf("- Pending pods: %d\n", resp.GetPendingPods()))
	b.WriteString(fmt.Sprintf("- Failed pods: %d\n", resp.GetFailedPods()))
	b.WriteString(fmt.Sprintf("- Succeeded pods: %d\n", resp.GetSucceededPods()))
	b.WriteString(fmt.Sprintf("- Unknown pods: %d\n", resp.GetUnknownPods()))

	if resp.GetUsage() != nil {
		b.WriteString(fmt.Sprintf("- Usage: CPU=%s, Memory=%s\n", resp.GetUsage().GetCpu(), resp.GetUsage().GetMemory()))
	}
	if resp.GetRequests() != nil {
		b.WriteString(fmt.Sprintf("- Requests: CPU=%s, Memory=%s\n", resp.GetRequests().GetCpu(), resp.GetRequests().GetMemory()))
	}
	if resp.GetLimits() != nil {
		b.WriteString(fmt.Sprintf("- Limits: CPU=%s, Memory=%s\n", resp.GetLimits().GetCpu(), resp.GetLimits().GetMemory()))
	}
}

func writeWorkloadsByHealth(b *strings.Builder, data any) {
	resp, ok := data.(*pb.WorkloadListResponse)
	if !ok {
		writeRawData(b, "Workloads By Health", data)
		return
	}

	workloads := resp.GetWorkloads()

	b.WriteString("Workloads By Health:\n")
	b.WriteString(fmt.Sprintf("- Total workloads returned: %d\n", len(workloads)))

	limit := min(len(workloads), maxWorkloadsInContext)
	for i := 0; i < limit; i++ {
		w := workloads[i]

		b.WriteString(fmt.Sprintf(
			"- %s/%s health=%s desired=%d ready=%d available=%d updated=%d active=%d succeeded=%d failed=%d age=%s\n",
			w.GetNamespace(),
			w.GetName(),
			w.GetHealth().String(),
			w.GetDesiredReplicas(),
			w.GetReadyReplicas(),
			w.GetAvailableReplicas(),
			w.GetUpdatedReplicas(),
			w.GetActive(),
			w.GetSucceeded(),
			w.GetFailed(),
			w.GetAge(),
		))

		if len(w.GetConditions()) > 0 {
			b.WriteString(fmt.Sprintf("  conditions: %s\n", strings.Join(w.GetConditions(), ", ")))
		}
	}

	if len(workloads) > limit {
		b.WriteString(fmt.Sprintf("- %d more workloads omitted\n", len(workloads)-limit))
	}
}

func writePods(b *strings.Builder, data any) {
	resp, ok := data.(*pb.PodListResponse)
	if !ok {
		writeRawData(b, "Pods", data)
		return
	}

	pods := resp.GetPods()

	statusCount := map[string]int{}
	for _, pod := range pods {
		statusCount[pod.GetStatus()]++
	}

	b.WriteString("Pods:\n")
	b.WriteString(fmt.Sprintf("- Total pods returned: %d\n", len(pods)))
	b.WriteString("- Status counts:\n")

	for _, status := range sortedKeys(statusCount) {
		b.WriteString(fmt.Sprintf("  - %s: %d\n", status, statusCount[status]))
	}

	b.WriteString("- Not-running pods:\n")

	written := 0
	for _, pod := range pods {
		if strings.EqualFold(pod.GetStatus(), "Running") {
			continue
		}

		if written >= maxPodsInContext {
			break
		}

		b.WriteString(fmt.Sprintf(
			"  - %s/%s status=%s node=%s\n",
			pod.GetNamespace(),
			pod.GetName(),
			pod.GetStatus(),
			pod.GetNodeName(),
		))
		written++
	}

	if written == 0 {
		b.WriteString("  - None found in returned data\n")
	}
}

func writePodStats(b *strings.Builder, data any) {
	stats, ok := data.([]*pb.PodStatsResponse)
	if !ok {
		writeRawData(b, "Pod Stats", data)
		return
	}

	b.WriteString("Pod Stats:\n")
	b.WriteString(fmt.Sprintf("- Total pod stats returned: %d\n", len(stats)))

	limit := min(len(stats), maxPodStatsInContext)
	for i := 0; i < limit; i++ {
		item := stats[i]
		b.WriteString(fmt.Sprintf(
			"- %s/%s CPU=%s Memory=%s\n",
			item.GetNamespace(),
			item.GetName(),
			item.GetCpu(),
			item.GetMemory(),
		))
	}

	if len(stats) > limit {
		b.WriteString(fmt.Sprintf("- %d more pod stats omitted\n", len(stats)-limit))
	}
}

func writeRecentEvents(b *strings.Builder, data any) {
	events, ok := data.([]*pb.EventResponse)
	if !ok {
		writeRawData(b, "Recent Events", data)
		return
	}

	b.WriteString("Recent Events:\n")
	b.WriteString(fmt.Sprintf("- Total events returned: %d\n", len(events)))

	limit := min(len(events), maxEventsInContext)
	for i := 0; i < limit; i++ {
		event := events[i]
		b.WriteString(fmt.Sprintf(
			"- Type=%s Reason=%s Object=%s Message=%s\n",
			event.GetType(),
			event.GetReason(),
			event.GetInvolvedObject(),
			shorten(event.GetMessage(), 300),
		))
	}

	if len(events) > limit {
		b.WriteString(fmt.Sprintf("- %d more events omitted\n", len(events)-limit))
	}
}

func writeNodes(b *strings.Builder, data any) {
	resp, ok := data.(*pb.NodeListResponse)
	if !ok {
		writeRawData(b, "Nodes", data)
		return
	}

	nodes := resp.GetNodes()

	b.WriteString("Nodes:\n")
	b.WriteString(fmt.Sprintf("- Total nodes: %d\n", len(nodes)))

	for _, node := range nodes {
		b.WriteString(fmt.Sprintf(
			"- %s internalIP=%s phase=%s age=%s\n",
			node.GetName(),
			node.GetInternalIp(),
			node.GetPhase(),
			node.GetAge(),
		))
	}
}

func writeNodeStats(b *strings.Builder, data any) {
	resp, ok := data.(*pb.NodeStatsResponse)
	if !ok {
		writeRawData(b, "Node Stats", data)
		return
	}

	b.WriteString("Node Stats:\n")
	b.WriteString(fmt.Sprintf("- Node: %s\n", resp.GetName()))
	b.WriteString(fmt.Sprintf("- CPU: %s\n", resp.GetCpu()))
	b.WriteString(fmt.Sprintf("- Memory: %s\n", resp.GetMemory()))
}

func writeNodeResourceAllocation(b *strings.Builder, data any) {
	resp, ok := data.(*pb.NodeResourceAllocationResponse)
	if !ok {
		writeRawData(b, "Node Resource Allocation", data)
		return
	}

	b.WriteString("Node Resource Allocation:\n")
	b.WriteString(fmt.Sprintf("- Node: %s\n", resp.GetNodeName()))
	b.WriteString(fmt.Sprintf("- Pod count: %d\n", resp.GetPodCount()))

	if resp.GetCapacity() != nil {
		b.WriteString(fmt.Sprintf("- Capacity: CPU=%s, Memory=%s\n", resp.GetCapacity().GetCpu(), resp.GetCapacity().GetMemory()))
	}
	if resp.GetAllocatable() != nil {
		b.WriteString(fmt.Sprintf("- Allocatable: CPU=%s, Memory=%s\n", resp.GetAllocatable().GetCpu(), resp.GetAllocatable().GetMemory()))
	}
	if resp.GetRequests() != nil {
		b.WriteString(fmt.Sprintf("- Requests: CPU=%s, Memory=%s\n", resp.GetRequests().GetCpu(), resp.GetRequests().GetMemory()))
	}
	if resp.GetLimits() != nil {
		b.WriteString(fmt.Sprintf("- Limits: CPU=%s, Memory=%s\n", resp.GetLimits().GetCpu(), resp.GetLimits().GetMemory()))
	}
	if resp.GetUsage() != nil {
		b.WriteString(fmt.Sprintf("- Usage: CPU=%s, Memory=%s\n", resp.GetUsage().GetCpu(), resp.GetUsage().GetMemory()))
	}
}

func writeFallback(b *strings.Builder, result ToolResult) {
	writeRawData(b, result.Name, result.Data)
}

func writeRawData(b *strings.Builder, title string, data any) {
	b.WriteString(title)
	b.WriteString(":\n")

	raw, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		b.WriteString("- unable to serialize result\n")
		return
	}

	text := string(raw)
	if len(text) > 2000 {
		text = text[:2000] + "\n[truncated]\n"
	}

	b.WriteString(text)
	b.WriteString("\n")
}

func sortedKeys(m map[string]int) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func shorten(s string, max int) string {
	s = strings.TrimSpace(s)
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}
