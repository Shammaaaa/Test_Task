package domain

type Note struct {
	ID     int64  `json:"id"`
	UserID int64  `json:"user_id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}
