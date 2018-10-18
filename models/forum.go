package models

// Информация о форуме.
type Forum struct {
	// Название форума.
	Title string `json:"title"`
	// Nickname пользователя, который отвечает за форум.
	User string `json:"user"`
	// Человекопонятный URL
	Slug string `json:"slug"`
	// Общее кол-во сообщений в данном форуме.
	Posts float32 `json:"posts,omitempty"`
	// Общее кол-во ветвей обсуждения в данном форуме.
	Threads float32 `json:"threads,omitempty"`
}
