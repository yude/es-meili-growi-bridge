package elasticsearch

type RootResponse struct {
	Name        string `json:"name"`
	ClusterName string `json:"cluster_name"`
	ClusterUUID string `json:"cluster_uuid"`
	Version     VersionInfo `json:"version"`
	TagLine     string `json:"tagline"`
}

type VersionInfo struct {
	Number        string `json:"number"`
	BuildFlavor   string `json:"build_flavor"`
	BuildType     string `json:"build_type"`
	BuildHash     string `json:"build_hash"`
	BuildDate     string `json:"build_date"`
	BuildSnapshot bool   `json:"build_snapshot"`
	LuceneVersion string `json:"lucene_version"`
	MinimumWireCompatibilityVersion string `json:"minimum_wire_compatibility_version"`
	MinimumIndexCompatibilityVersion string `json:"minimum_index_compatibility_version"`
}

func DefaultRootResponse() *RootResponse {
	return &RootResponse{
		Name:        "es-meili-growi-bridge",
		ClusterName: "es-meili-growi-bridge",
		ClusterUUID: "BDi0gxodQSeWOKftHUK6Lg",
		Version: VersionInfo{
			Number:                        "8.17.0",
			BuildFlavor:                   "default",
			BuildType:                     "docker",
			BuildHash:                     "f11129f23bc31112e4f0a13e6e0d8b7edc0a21a2",
			BuildDate:                     "2025-12-10T10:10:00.000000000Z",
			BuildSnapshot:                 false,
			LuceneVersion:                 "9.12.0",
			MinimumWireCompatibilityVersion: "7.17.0",
			MinimumIndexCompatibilityVersion: "7.0.0",
		},
		TagLine: "You Know, for Search (powered by Meilisearch)",
	}
}
