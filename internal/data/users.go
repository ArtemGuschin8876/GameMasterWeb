package data

type Users struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	City     string `json:"city"`
	About    string `json:"about"`
	Image    string `json:"image"`
}
