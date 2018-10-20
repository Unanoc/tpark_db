package errors

import (
	"encoding/json"
	"log"
)

//easyjson:json
type Error struct {
	Message string `json:"message,omitempty"`
}

func (r *Error) Error() string {
	errorBytes, err := json.Marshal(r)
	if err != nil {
		log.Println(err)
	}
	return string(errorBytes)
}

func New(msg string) error {
	return &Error{Message: msg}
}

var ForumIsExist = New("Forum was created earlier")
var UserNotFound = New("User not found")
var ForumNotFound = New("Forum not found")
var ThreadIsExist = New("Thread was created earlier")
var UserIsExist = New("User was created earlier")
var UserUpdateConflict = New("User not updated")
var ThreadNotFound = New("Thread not found")
var NoPostsForCreate = New("Not posts for create")
var NoThreadParent = New("No parent for thread")
var PostNotFound = New("Post not found")
