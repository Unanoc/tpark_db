package models

import (
	"tpark_db/database"
	"tpark_db/errors"
)

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

func (f *Forum) CreateForum() (*Forum, error) {
	tx := database.StartTransaction()
	defer tx.Rollback()

	isForumExist := Forum{}
	if err := tx.QueryRow(`
		SELECT "title", "user", "slug", "posts", "threads"
		WHERE "slug" = $1 AND "user" = $2 AND "title" = $3`,
		f.Slug, f.User, f.Title).Scan(&isForumExist.Posts,
		&isForumExist.Slug, &isForumExist.Threads,
		&isForumExist.Title, &isForumExist.User); err == nil {
		return &isForumExist, errors.ForumIsExist
	}

	rows := tx.QueryRow(`
		INSERT
		INTO forums ("slug", "title", "user")
		VALUES ($1, $2, (SELECT nickname FROM users WHERE nickname = $3)) // tck
		RETURNING "user"`,
		f.Slug, f.Title, f.User)
	if err := rows.Scan(&f.User); err != nil {
		return nil, errors.UserNotFound
	}

	database.CommitTransaction(tx)
	return nil, nil
}
