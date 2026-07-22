package elasticsearch

type ShardStatistics struct {
	Total      int `json:"total"`
	Successful int `json:"successful"`
	Skipped    int `json:"skipped"`
	Failed     int `json:"failed"`
}

func DefaultShards() ShardStatistics {
	return ShardStatistics{Total: 1, Successful: 1, Skipped: 0, Failed: 0}
}
