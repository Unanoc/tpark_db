package helpers

import (
	"strconv"
	"time"
	"tpark_db/database"
	"tpark_db/errors"
	"tpark_db/models"

	"github.com/jackc/pgx"
)

func ThreadCreateHelper(posts *models.Posts, slugOrId string) (*models.Posts, error) {
	if len(*posts) == 0 {
		return nil, errors.NoPostsForCreate
	}

	var threadByID *models.Thread
	var err error
	if IsNumber(slugOrId) {
		id, _ := strconv.Atoi(slugOrId)
		threadByID, err = GetThreadById(id)
		if err != nil {
			return nil, err
		}
	} else {
		threadByID, err = GetThreadBySlug(slugOrId)
		if err != nil {
			return nil, err
		}
	}

	tx := database.StartTransaction()
	defer tx.Rollback()

	created := time.Now()
	insertedPosts := models.Posts{}
	for _, post := range *posts {
		var rows *pgx.Row

		if post.Parent != 0 {
			// TODO
		} else {
			rows = tx.QueryRow(`
				INSERT
				INTO posts (author, created, message, thread, parent, forum)
				VALUES ($1, $2, $3, $4, $5, $6) 
				RETURNING author, created, forum, id, message, parent, thread`,
				post.Author,
				created,
				post.Message,
				threadByID.Id,
				post.Parent,
				threadByID.Forum,
			)

			insertedPost := models.Post{}
			err := rows.Scan(
				&insertedPost.Author,
				&insertedPost.Created,
				&insertedPost.Forum,
				&insertedPost.Id,
				&insertedPost.Message,
				&insertedPost.Parent,
				&insertedPost.Thread,
			)
			if err != nil {
				return nil, err
			}
			insertedPosts = append(insertedPosts, &insertedPost)
		}

	}

	database.CommitTransaction(tx)
	return &insertedPosts, nil
}
