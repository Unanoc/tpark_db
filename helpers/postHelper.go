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
		return nil, err
	}

	postFull := models.PostFull{}

	for _, typeObject := range related {
		switch typeObject {
		default:
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
			fmt.Println(err)
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
		SELECT id, author, message, forum, thread, created
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
	)

	if err != nil {
		return nil, errors.PostNotFound
	}

	database.CommitTransaction(tx)
	return &post, nil
}
