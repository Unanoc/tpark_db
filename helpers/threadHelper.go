package helpers

import (
	"fmt"
	"strconv"
	"time"
	"tpark_db/database"
	"tpark_db/errors"
	"tpark_db/models"
)

func ThreadCreateHelper(posts *models.Posts, slug string) error {
	tx := database.StartTransaction()
	defer tx.Rollback()

	if len(*posts) == 0 {
		return errors.NoPostsForCreate
	}

	batch := tx.BeginBatch()
	defer batch.Close()

	// checking if thread exist is database
	id, err := strconv.Atoi(slug) // checking is "slug" or "id"
	var forum string
	if err == nil {
		rows := tx.QueryRow(` 
			SELECT id, forum
			FROM threads
			WHERE id = $1
		`, id)
		if err = rows.Scan(&id, &forum); err != nil {
			return errors.ThreadNotFound
		}
	} else {
		rows := tx.QueryRow(`
			SELECT id, forum
			FROM threads
			WHERE slug = $1
		`, slug)
		if err = rows.Scan(&id, &forum); err != nil {
			return errors.ThreadNotFound
		}
	}

	created := time.Now()

	//TODOs
	for _, post := range *posts {

		rows := tx.QueryRow(`
			INSERT
			INTO posts (author, created, message, thread, parent)
			VALUES ($1, $2, $3, $4, $5) 
			RETURNING author, created, forum, id, isEdited, message, parent, thread`,
			&post.Author,
			created,
			&post.Message,
			&post.Thread,
			&post.Parent)

		err := rows.Scan(&post.Author, &post.Created, &post.Forum, &post.Id, &post.IsEdited, &post.Message, &post.Parent, &post.Thread)
		if err != nil {
			// switch err.(pgx.PgError).Code {
			// case "23505":
			// 	return errors.ThreadIsExist
			// case "23502":
			// 	return errors.UserNotFound
			// default:
			// 	return err
			// }

		}

	}

	database.CommitTransaction(tx)
	return nil
}
