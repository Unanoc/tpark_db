package helpers

import (
	"fmt"
	"strconv"
	"tpark_db/database"
	"tpark_db/errors"
	"tpark_db/models"
)

func PostFullHelper(id string, related []string) (*models.PostFull, error) {
	postID, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	postFull := models.PostFull{}

	for _, typeObject := range related {
		switch typeObject {
		case "thread":
			threadID := strconv.Itoa(postFull.Post.Thread)
			postFull.Thread, err = GetThreadBySlugOrId(threadID)
		case "forum":
			forumSlug := postFull.Post.Forum
			postFull.Forum, err = ForumGetBySlug(forumSlug)
		case "user":
			userNickname := postFull.Post.Author
			postFull.Author, err = UserGetOneHelper(userNickname)
		default:
			postFull.Post, err = PostGetOneById(postID)
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

	rows := tx.QueryRow(` 
		SELECT id, author, message, forum, thread, created, isEdited
		FROM posts
		WHERE id = $1`,
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
		return nil, errors.PostNotFound
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

	if len(postUpdate.Message) != 0 && (post.Message != postUpdate.Message) {
		post.Message = postUpdate.Message
		post.IsEdited = true
	} else {
		return post, nil
	}

	tx := database.StartTransaction()
	defer tx.Rollback()

	_, err = tx.Exec(`
		UPDATE posts
		SET message = coalesce(nullif($2, ''), message),
			isEdited = TRUE
		WHERE id = $1`,
		postID, &postUpdate.Message)

	if err != nil {
		return nil, err
	}

	database.CommitTransaction(tx)
	return post, nil
}
