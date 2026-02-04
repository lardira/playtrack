package domain

type ResponseMessage struct {
	Body struct {
		Message string `json:"message"`
	}
}

type ResponseItems[T any] struct {
	Body struct {
		Items []T `json:"items"`
	}
}

type ResponseItem[T any] struct {
	Body struct {
		Item *T `json:"item"`
	}
}

type ResponseID[T any] struct {
	Body struct {
		ID T `json:"id"`
	}
}
