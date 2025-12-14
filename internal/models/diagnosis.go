package models

type HealthReport struct {
	Method      string   `json:"method"`
	Status      string   `json:"status"`
	HealthScore int      `json:"health_score"`
	Summary     string   `json:"summary"`
	Problems    []string `json:"problems"`
	Metrics     KeyStats `json:"metrics"`
}

type KeyStats struct {
	Temperature     int     `json:"temperature_c"`
	LifeRemaining   int     `json:"life_remaining_percent"`
	DataWrittenTB   float64 `json:"data_written_tb"`
	PowerOnHours    int     `json:"power_on_hours"`
	MediaErrors     int     `json:"media_errors"`
	UnsafeShutdowns int     `json:"unsafe_shutdowns"`
}
