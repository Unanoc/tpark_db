package helpers

import (
	"tpark_db/database"
	"tpark_db/errors"
	"tpark_db/models"
)

func ForumCreateHelper(f *models.Forum) (*models.Forum, error) {
	tx := database.StartTransaction()
	defer tx.Rollback()

	isForumExist := models.Forum{}
	if err := tx.QueryRow(`
		SELECT "title", "user", "slug", "posts", "threads" 
		FROM forum
		WHERE "slug" = $1 AND "user" = $2 AND "title" = $3`,
		f.Slug, f.User, f.Title).Scan(&isForumExist.Posts,
		&isForumExist.Slug, &isForumExist.Threads,
		&isForumExist.Title, &isForumExist.User); err == nil {
		return &isForumExist, errors.ForumIsExist
	}

	rows := tx.QueryRow(`
		INSERT
		INTO forums ("slug", "title", "user")
		VALUES ($1, $2, (SELECT nickname FROM users WHERE nickname = $3)) 
		RETURNING "user"`,
		f.Slug, f.Title, f.User)
	if err := rows.Scan(&f.User); err != nil {
		return nil, errors.UserNotFound
	}

	database.CommitTransaction(tx)
	return nil, nil
}
