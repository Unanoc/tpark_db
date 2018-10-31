package helpers

import (
	"strconv"
	"tpark_db/database"
	"tpark_db/errors"
	"tpark_db/models"
)

func PostFullHelper(id string, related []string) (*models.PostFull, error) {
	postID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	postFull := models.PostFull{}

	for _, typeObject := range related {
		switch typeObject {
		case "post":
			postFull.Post, err = PostGetOneById(postID)
		case "thread":
			threadID := strconv.Itoa(postFull.Post.Thread)
			postFull.Thread, err = GetThreadBySlugOrId(threadID)
		case "forum":
			forumSlug := postFull.Post.Forum
			postFull.Forum, err = ForumGetBySlug(forumSlug)
		case "user":
			userNickname := postFull.Post.Author
			postFull.Author, err = UserGetOneHelper(userNickname)
		}

		if err != nil {
			return nil, err
		}
	}

	return &postFull, nil
}

func PostGetOneById(id int) (*models.Post, error) {
	tx := database.StartTransaction()
	defer tx.Rollback()

	post := models.Post{}

	sql := "SELECT id, author, \"message\", forum, thread, created, \"isEdited\" FROM posts WHERE id = $1"
	rows := tx.QueryRow(sql,
		id)

	err := rows.Scan(
		&post.Id,
		&post.Author,
		&post.Message,
		&post.Forum,
		&post.Thread,
		&post.Created,
		&post.IsEdited,
	)

	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, errors.PostNotFound
		}
		return nil, err
	}

	database.CommitTransaction(tx)
	return &post, nil
}

func PostUpdateHelper(postUpdate *models.PostUpdate, postID string) (*models.Post, error) {
	id, err := strconv.Atoi(postID)
	if err != nil {
		return nil, err
	}

	post, err := PostGetOneById(id)
	if err != nil {
		return nil, errors.PostNotFound
	}

	if len(postUpdate.Message) == 0 {
		return post, nil
	}

	tx := database.StartTransaction()
	defer tx.Rollback()

	sql := "UPDATE posts SET \"message\" = COALESCE($2, \"message\"), \"isEdited\" = ($2 IS NOT NULL AND $2 <> \"message\") WHERE id = $1 RETURNING author::text, created, forum, \"isEdited\", thread, \"message\""
	rows := tx.QueryRow(sql,
		postID, &postUpdate.Message)

	err = rows.Scan(
		&post.Author,
		&post.Created,
		&post.Forum,
		&post.IsEdited,
		&post.Thread,
		&post.Message,
	)

	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, errors.PostNotFound
		}
		return nil, err
	}

	database.CommitTransaction(tx)
	return post, nil
}
