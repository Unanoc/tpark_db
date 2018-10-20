package models

//easyjson:json
type User struct {
	// Имя пользователя (уникальное поле). Данное поле допускает только латиницу, цифры и знак подчеркивания. Сравнение имени регистронезависимо.
	Nickname string `json:"nickname,omitempty"`
	// Полное имя пользователя.
	Fullname string `json:"fullname"`
	// Описание пользователя.
	About string `json:"about,omitempty"`
	// Почтовый адрес пользователя (уникальное поле).
	Email string `json:"email"`
}

//easyjson:json
type UserUpdate struct {
	// Полное имя пользователя.
	Fullname string `json:"fullname,omitempty"`
	// Описание пользователя.
	About string `json:"about,omitempty"`
	// Почтовый адрес пользователя (уникальное поле).
	Email string `json:"email,omitempty"`
}

//easyjson:json
type Users []*User
