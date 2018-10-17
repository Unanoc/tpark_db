package models

import (
	"fmt"
	"log"
	"tpark_db/database"
	"tpark_db/errors"
)

// Информация о пользователе.
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

// Информация о пользователе.
type UserUpdate struct {
	// Полное имя пользователя.
	Fullname string `json:"fullname,omitempty"`
	// Описание пользователя.
	About string `json:"about,omitempty"`
	// Почтовый адрес пользователя (уникальное поле).
	Email string `json:"email,omitempty"`
}

type Users []*User

func (u *User) CreateUser() (Users, error) {
	tx := database.StartTransaction()
	defer tx.Rollback()

	rows, err := tx.Exec(`
		INSERT
		INTO users ("nickname", "fullname", "about", "email")
		VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`,
		&u.Nickname, &u.Fullname, &u.About, &u.Email)

	if err != nil {
		log.Println(err)
		return nil, err
	}
	fmt.Println("here")
	if rows.RowsAffected() == 0 { // if it returns 0 - user existed, else user was created
		users := Users{}
		queryRows, err := tx.Query(`
			SELECT "nickname", "fullname", "about", "email"
			FROM users
			WHERE "email" = $1 OR "nickname" = $2`,
			&u.Email, &u.Nickname)

		defer queryRows.Close()

		if err != nil {
			log.Println(err)
			return nil, err
		}

		for queryRows.Next() {
			user := User{}
			queryRows.Scan(&user.Nickname, &user.Fullname,
				&user.About, &user.Email)
			users = append(users, &user)
		}
		return users, errors.UserIsExist
	}

	database.CommitTransaction(tx)
	return nil, nil
}
