package models

//easyjson:json
type Vote struct {
	Nickname string  `json:"nickname"`
	Voice    float32 `json:"voice"`
}
