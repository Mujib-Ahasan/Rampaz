package ai

type ToolSchema struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Parameters  map[string]any `json:"parameters"`
}

func AvailableTools() []ToolSchema {
	return []ToolSchema{
		{
			Name:        "get_cluster_overview",
			Description: "Get high-level cluster-wide resource counts such as nodes, namespaces, pods, workloads, services, PVCs, and network policies.",
			Parameters:  objectSchema(nil, nil),
		},
		{
			Name:        "get_namespace_summary",
			Description: "Get resource counts for a specific Kubernetes namespace, including pods, workloads, services, PVCs, and network policies.",
			Parameters: objectSchema(map[string]any{
				"namespace": stringParam("Kubernetes namespace name."),
			}, []string{"namespace"}),
		},
		{
			Name:        "get_namespace_metrics",
			Description: "Get namespace-level pod status counts and resource metrics such as CPU and memory usage, requests, and limits.",
			Parameters: objectSchema(map[string]any{
				"namespace": stringParam("Kubernetes namespace name."),
			}, []string{"namespace"}),
		},
		{
			Name:        "get_workloads_by_health",
			Description: "Get workloads filtered by health status, optionally scoped to a namespace. Useful for finding degraded or unhealthy workloads.",
			Parameters: objectSchema(map[string]any{
				"health": map[string]any{
					"type":        "string",
					"description": "Workload health status.",
					"enum":        []string{"HEALTHY", "DEGRADED", "UNHEALTHY"},
				},
				"namespace": stringParam("Optional Kubernetes namespace name."),
			}, []string{"health"}),
		},
		{
			Name:        "list_pods",
			Description: "List pods, optionally scoped to a namespace and label selector. Useful for checking pod status, node placement, and failed/pending pods.",
			Parameters: objectSchema(map[string]any{
				"namespace":      stringParam("Optional Kubernetes namespace name. Empty means all namespaces if supported by backend."),
				"label_selector": stringParam("Optional Kubernetes label selector, for example app=nginx."),
			}, nil),
		},
		{
			Name:        "get_recent_events",
			Description: "Get recent Kubernetes events, optionally scoped to a namespace. Useful for debugging scheduling, image pull, crash, and readiness issues.",
			Parameters: objectSchema(map[string]any{
				"namespace": stringParam("Optional Kubernetes namespace name."),
				"limit": map[string]any{
					"type":        "integer",
					"description": "Maximum number of events to return.",
					"minimum":     1,
					"maximum":     20,
					"default":     20,
				},
			}, nil),
		},
		{
			Name:        "list_nodes",
			Description: "List Kubernetes nodes with basic metadata such as name, internal IP, phase, and age.",
			Parameters:  objectSchema(nil, nil),
		},
		{
			Name:        "get_node_stats",
			Description: "Get current CPU and memory usage for a specific node.",
			Parameters: objectSchema(map[string]any{
				"nodeName": stringParam("Kubernetes node name."),
			}, []string{"nodeName"}),
		},
		{
			Name:        "get_node_resource_allocation",
			Description: "Get node capacity, allocatable resources, requests, limits, usage, and pod count for a specific node.",
			Parameters: objectSchema(map[string]any{
				"nodeName": stringParam("Kubernetes node name."),
			}, []string{"nodeName"}),
		},
		{
			Name:        "get_pod_stats",
			Description: "Get current CPU and memory usage for pods, optionally scoped to a namespace. Results are bounded by limit.",
			Parameters: objectSchema(map[string]any{
				"namespace": stringParam("Optional Kubernetes namespace name. Empty means all namespaces if supported by backend."),
				"limit": map[string]any{
					"type":        "integer",
					"description": "Maximum number of pod stats items to return.",
					"minimum":     1,
					"maximum":     50,
					"default":     20,
				},
			}, nil),
		},
		{
			Name:        "get_pod_logs",
			Description: "Fetch recent error and warning logs for a specific pod. Use this to debug pod failures, crashes, or runtime issues.",
			Parameters: objectSchema(map[string]any{
				"namespace": stringParam("Must give a namespace name in which the pod belongs"),
				"pod_name":  stringParam("Pod name for which logs are needed"),
			}, []string{"namespace", "pod_name"}),
		},
	}
}

func objectSchema(properties map[string]any, required []string) map[string]any {
	if properties == nil {
		properties = map[string]any{}
	}

	schema := map[string]any{
		"type":                 "object",
		"properties":           properties,
		"additionalProperties": false,
	}

	if len(required) > 0 {
		schema["required"] = required
	}

	return schema
}

func stringParam(description string) map[string]any {
	return map[string]any{
		"type":        "string",
		"description": description,
	}
}
