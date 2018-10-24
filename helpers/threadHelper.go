package helpers

import (
	"fmt"
	"strconv"
	"time"
	"tpark_db/database"
	"tpark_db/errors"
	"tpark_db/models"

	"github.com/jackc/pgx"
)

func GetThreadBySlug(slug string) (*string, error) {
	tx := database.StartTransaction()
	defer tx.Rollback()

	rows := tx.QueryRow(` 
		SELECT slug
		FROM threads
		WHERE slug = $1
	`, slug)

	var result string
	err := rows.Scan(&result)
	if err != nil {
		return nil, errors.ThreadNotFound
	}

	database.CommitTransaction(tx)
	return &result, nil
}

func GetThreadById(id int) (*string, error) {
	tx := database.StartTransaction()
	defer tx.Rollback()

	rows := tx.QueryRow(` 
		SELECT slug
		FROM threads
		WHERE id = $1
	`, id)

	var result string
	err := rows.Scan(&result)
	if err != nil {
		return nil, errors.ThreadNotFound
	}

	database.CommitTransaction(tx)
	return &result, nil
}

func ThreadCreateHelper(posts *models.Posts, slugOrId string) (*models.Posts, error) {
	if len(*posts) == 0 {
		return nil, errors.NoPostsForCreate
	}

	var slugThread *string
	var err error
	if IsNumber(slugOrId) {
		id, _ := strconv.Atoi(slugOrId)
		slugThread, err = GetThreadById(id)
		if err != nil {
			return nil, err
		}
	} else {
		slugThread, err = GetThreadBySlug(slugOrId)
		if err != nil {
			return nil, err
		}
	}

	fmt.Println(slugThread)

	tx := database.StartTransaction()
	defer tx.Rollback()

	created := time.Now()
	insertedPosts := models.Posts{}
	for _, post := range *posts {
		var rows *pgx.Row

		if post.Parent != 0 {
			rows = tx.QueryRow(`
				INSERT
				INTO posts (author, created, message, thread, parent, forum, id, isEdited)
				VALUES ($1, $2, $3, $4, $5) 
				RETURNING author, created, forum, id, isEdited, message, parent, thread`,
				&post.Author,
				&created,
				&post.Message,
				&slugThread,
				&post.Parent,
			)

			insertedPost := models.Post{}
			err := rows.Scan(
				&insertedPost.Author,
				&insertedPost.Created,
				&insertedPost.Forum,
				&insertedPost.Id,
				&insertedPost.IsEdited,
				&insertedPost.Message,
				&insertedPost.Parent,
				&insertedPost.Thread,
			)
			if err != nil {
				return nil, errors.PostParentNotFound
			}

		} else {
			rows = tx.QueryRow(`
				INSERT
				INTO posts (author, created, message, thread, parent, forum, id, isEdited)
				VALUES ($1, $2, $3, $4, $5) 
				RETURNING author, created, forum, id, isEdited, message, parent, thread`,
				&post.Author,
				&created,
				&post.Message,
				&slugThread,
				&post.Parent,
			)

			insertedPost := models.Post{}
			err := rows.Scan(
				&insertedPost.Author,
				&insertedPost.Created,
				&insertedPost.Forum,
				&insertedPost.Id,
				&insertedPost.IsEdited,
				&insertedPost.Message,
				&insertedPost.Parent,
				&insertedPost.Thread,
			)
			if err != nil {
				return nil, err
			}
		}

	}

	database.CommitTransaction(tx)
	return &insertedPosts, nil
}
