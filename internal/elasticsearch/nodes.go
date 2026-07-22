package elasticsearch

type NodesInfoResponse struct {
	ClusterName string             `json:"cluster_name"`
	Nodes       map[string]NodeInfo `json:"nodes"`
}

type NodeInfo struct {
	Name             string          `json:"name"`
	TransportAddress string          `json:"transport_address"`
	Host             string          `json:"host"`
	IP               string          `json:"ip"`
	Version          string          `json:"version"`
	BuildFlavor      string          `json:"build_flavor"`
	BuildType        string          `json:"build_type"`
	HTTPAddress      string          `json:"http_address"`
	Roles            []string        `json:"roles"`
	Plugins          []PluginInfo    `json:"plugins"`
	Modules          []PluginInfo    `json:"modules"`
}

type PluginInfo struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
}

func DefaultNodesInfo() *NodesInfoResponse {
	return &NodesInfoResponse{
		ClusterName: "es-meili-growi-bridge",
		Nodes: map[string]NodeInfo{
			"meilisearch-node-1": {
				Name:             "meilisearch-0",
				TransportAddress: "0.0.0.0:9300",
				Host:             "0.0.0.0",
				IP:               "127.0.0.1",
				Version:          "8.17.0",
				BuildFlavor:      "default",
				BuildType:        "docker",
				HTTPAddress:      "0.0.0.0:9200",
				Roles:            []string{"master", "data", "ingest"},
				Plugins:          []PluginInfo{},
				Modules:          []PluginInfo{},
			},
		},
	}
}
