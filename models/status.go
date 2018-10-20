package models

//easyjson:json
type Status struct {
	User   float32 `json:"user"`
	Forum  float32 `json:"forum"`
	Thread float32 `json:"thread"`
	Post   float32 `json:"post"`
}
