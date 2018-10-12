package models

import (
	"time"
)

// Ветка обсуждения на форуме.
type Thread struct {
	// Идентификатор ветки обсуждения.
	Id float32 `json:"id,omitempty"`
	// Заголовок ветки обсуждения.
	Title string `json:"title"`
	// Пользователь, создавший данную тему.
	Author string `json:"author"`
	// Форум, в котором расположена данная ветка обсуждения.
	Forum string `json:"forum,omitempty"`
	// Описание ветки обсуждения.
	Message string `json:"message"`
	// Кол-во голосов непосредственно за данное сообщение форума.
	Votes float32 `json:"votes,omitempty"`
	// Человекопонятный URL 
	Slug string `json:"slug,omitempty"`
	// Дата создания ветки на форуме.
	Created time.Time `json:"created,omitempty"`
}

// Сообщение для обновления ветки обсуждения на форуме. Пустые параметры остаются без изменений.
type ThreadUpdate struct {
	// Заголовок ветки обсуждения.
	Title string `json:"title,omitempty"`
	// Описание ветки обсуждения.
	Message string `json:"message,omitempty"`
}

type Threads []*Thread

