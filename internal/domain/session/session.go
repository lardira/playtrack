package session

type Session struct {
	ID    string `json:"id" format:"uuid"`
	Title string `json:"title"`
}
