package domain

type MessageResponse struct {
	Body struct {
		Message string `json:"message"`
	}
}
