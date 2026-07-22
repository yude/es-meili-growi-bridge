package elasticsearch

type ClusterHealthResponse struct {
	ClusterName                 string  `json:"cluster_name"`
	Status                      string  `json:"status"`
	TimedOut                    bool    `json:"timed_out"`
	NumberOfNodes               int     `json:"number_of_nodes"`
	NumberOfDataNodes           int     `json:"number_of_data_nodes"`
	ActivePrimaryShards         int     `json:"active_primary_shards"`
	ActiveShards                int     `json:"active_shards"`
	RelocatingShards            int     `json:"relocating_shards"`
	InitializingShards          int     `json:"initializing_shards"`
	UnassignedShards            int     `json:"unassigned_shards"`
	DelayedUnassignedShards     int     `json:"delayed_unassigned_shards"`
	NumberOfPendingTasks        int     `json:"number_of_pending_tasks"`
	NumberOfInFlightFetch       int     `json:"number_of_in_flight_fetch"`
	TaskMaxWaitingInQueueMillis int     `json:"task_max_waiting_in_queue_millis"`
	ActiveShardsPercentAsNumber float64 `json:"active_shards_percent_as_number"`
}

func DefaultClusterHealth() *ClusterHealthResponse {
	return &ClusterHealthResponse{
		ClusterName:                 "es-meili-growi-bridge",
		Status:                      "green",
		TimedOut:                    false,
		NumberOfNodes:               1,
		NumberOfDataNodes:           1,
		ActivePrimaryShards:         0,
		ActiveShards:                0,
		RelocatingShards:            0,
		InitializingShards:          0,
		UnassignedShards:            0,
		DelayedUnassignedShards:     0,
		NumberOfPendingTasks:        0,
		NumberOfInFlightFetch:       0,
		TaskMaxWaitingInQueueMillis: 0,
		ActiveShardsPercentAsNumber: 100.0,
	}
}
