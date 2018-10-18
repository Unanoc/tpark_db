package helpers

import (
	"log"
	"tpark_db/database"
	"tpark_db/errors"
	"tpark_db/models"
)

func ForumCreateHelper(f *models.Forum) (*models.Forum, error) {
	tx := database.StartTransaction()
	defer tx.Rollback()

	rows := tx.QueryRow(`
		INSERT
		INTO forums (slug, title, "user")
		VALUES ($1, $2, (SELECT nickname FROM users WHERE nickname = $3)) 
		RETURNING "user"`,
		f.Slug, f.Title, f.User)

	if err := rows.Scan(&f.User); err != nil {
		sError := err.Error()
		if sError[len(sError)-2] == '5' { // determinatingn an error by last number of error msg: "duplicate key value violates unique constraint "users_email_key" (SQLSTATE 23505)". It is bad code...  like API
			forum, _ := existsForum(f.Slug)
			log.Println(err)
			return forum, errors.ForumIsExist
		}
		log.Println(err)
		return nil, errors.UserNotFound
	}

	database.CommitTransaction(tx)
	return f, nil
}

func existsForum(slug string) (*models.Forum, error) {
	tx := database.StartTransaction()
	defer tx.Rollback()

	forum := models.Forum{}
	if err := tx.QueryRow(`
		SELECT title, "user", slug, posts, threads
		FROM forum
		WHERE "slug" = $1`,
		slug).Scan(&forum.Title,
		&forum.User, &forum.Slug,
		&forum.Posts, &forum.Threads); err == nil {
		log.Println(err)
		return &forum, errors.ForumIsExist
	}

	database.CommitTransaction(tx)
	return &forum, nil
}
