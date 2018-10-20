package models

import (
	"time"
)

//easyjson:json
type Thread struct {
	Id      float32   `json:"id,omitempty"`
	Title   string    `json:"title"`
	Author  string    `json:"author"`
	Forum   string    `json:"forum,omitempty"`
	Message string    `json:"message"`
	Votes   float32   `json:"votes,omitempty"`
	Slug    string    `json:"slug,omitempty"`
	Created time.Time `json:"created,omitempty"`
}

//easyjson:json
type ThreadUpdate struct {
	Title   string `json:"title,omitempty"`
	Message string `json:"message,omitempty"`
}

//easyjson:json
type Threads []*Thread
