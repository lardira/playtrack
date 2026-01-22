package tech

type Status struct {
	DB     bool `json:"db"`
	Server bool `json:"server"`
}

type HealthResponse struct {
	Body struct {
		Status `json:"status"`
	}
}
