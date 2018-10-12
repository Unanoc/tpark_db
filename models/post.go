package models

import (
	"time"
)

// Сообщение внутри ветки обсуждения на форуме. 
type Post struct {
	// Идентификатор данного сообщения.
	Id float32 `json:"id,omitempty"`
	// Идентификатор родительского сообщения (0 - корневое сообщение обсуждения). 
	Parent float32 `json:"parent,omitempty"`
	// Автор, написавший данное сообщение.
	Author string `json:"author"`
	// Собственно сообщение форума.
	Message string `json:"message"`
	// Истина, если данное сообщение было изменено.
	IsEdited bool `json:"isEdited,omitempty"`
	// Идентификатор форума (slug) данного сообещния.
	Forum string `json:"forum,omitempty"`
	// Идентификатор ветви (id) обсуждения данного сообещния.
	Thread float32 `json:"thread,omitempty"`
	// Дата создания сообщения на форуме.
	Created time.Time `json:"created,omitempty"`
}

// Полная информация о сообщении, включая связанные объекты.
type PostFull struct {
	Post *Post `json:"post,omitempty"`
	Author *User `json:"author,omitempty"`
	Thread *Thread `json:"thread,omitempty"`
	Forum *Forum `json:"forum,omitempty"`
}

// Сообщение для обновления сообщения внутри ветки на форуме. Пустые параметры остаются без изменений.
type PostUpdate struct {
	// Собственно сообщение форума.
	Message string `json:"message,omitempty"`
}

type Posts []*Post
