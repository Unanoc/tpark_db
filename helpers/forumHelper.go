package helpers

import (
	"fmt"
	"tpark_db/database"
	"tpark_db/errors"
	"tpark_db/models"

	"github.com/jackc/pgx"
)

func ForumCreateHelper(f *models.Forum) (*models.Forum, error) {
	tx := database.StartTransaction()
	defer tx.Rollback()

	rows := tx.QueryRow(`
		INSERT
		INTO forums (slug, title, "user")
		VALUES ($1, $2, (SELECT nickname FROM users WHERE nickname = $3)) 
		RETURNING "user"`,
		f.Slug,
		f.Title,
		f.User)

	if err := rows.Scan(&f.User); err != nil {
		switch err.(pgx.PgError).Code {
		case "23505":
			forum, err := ForumGetBySlug(f.Slug)
			fmt.Println(err)
			return forum, errors.ForumIsExist
		case "23502":
			return nil, errors.UserNotFound
		default:
			return nil, err
		}
	}

	database.CommitTransaction(tx)
	return f, nil
}

func ForumGetBySlug(slug string) (*models.Forum, error) {
	tx := database.StartTransaction()
	defer tx.Rollback()

	forum := models.Forum{}

	err := tx.QueryRow(`
		SELECT slug, title, "user"
		FROM forums
		WHERE slug = $1`,
		slug).Scan(
		&forum.Slug,
		&forum.Title,
		&forum.User)

	if err != nil {
		return nil, err
	}

	database.CommitTransaction(tx)
	return &forum, nil
}
