package models

import (
	"time"
)

//easyjson:json
type Post struct {
	Id       int64     `json:"id,omitempty"`
	Parent   int64     `json:"parent,omitempty"`
	Author   string    `json:"author"`
	Message  string    `json:"message"`
	IsEdited bool      `json:"isEdited,omitempty"`
	Forum    string    `json:"forum,omitempty"`
	Thread   int       `json:"thread,omitempty"`
	Created  time.Time `json:"created,omitempty"`
	Path     []int64   `json:"-"`
}

//easyjson:json
type PostFull struct {
	Post   *Post   `json:"post,omitempty"`
	Author *User   `json:"author,omitempty"`
	Thread *Thread `json:"thread,omitempty"`
	Forum  *Forum  `json:"forum,omitempty"`
}

//easyjson:json
type PostUpdate struct {
	Message string `json:"message,omitempty"`
}

//easyjson:json
type Posts []*Post
