package helpers

import (
	"strconv"
	"tpark_db/database"
	"tpark_db/errors"
	"tpark_db/models"
)

// PostUpdateHelper updates post data.
func PostUpdateHelper(postUpdate *models.PostUpdate, postID string) (*models.Post, error) {
	id, err := strconv.Atoi(postID)
	if err != nil {
		return nil, err
	}

	post, err := PostGetOneByIDHelper(id)
	if err != nil {
		return nil, errors.PostNotFound
	}

	if len(postUpdate.Message) == 0 {
		return post, nil
	}

	rows := database.DB.Conn.QueryRow(sqlUpdatePost, postID, &postUpdate.Message)

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

	return post, nil
}

// PostFullHelper selects full data about post.
func PostFullHelper(id string, related []string) (*models.PostFull, error) {
	postID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	postFull := models.PostFull{}

	for _, typeObject := range related {
		switch typeObject {
		case "post":
			postFull.Post, err = PostGetOneByIDHelper(postID)
		case "thread":
			threadID := strconv.Itoa(postFull.Post.Thread)
			postFull.Thread, err = GetThreadBySlugOrID(threadID)
		case "forum":
			forumSlug := postFull.Post.Forum
			postFull.Forum, err = ForumGetBySlugHelper(forumSlug)
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

// PostGetOneByIDHelper selects post by id from table POSTS
func PostGetOneByIDHelper(id int) (*models.Post, error) {
	post := models.Post{}

	rows := database.DB.Conn.QueryRow(sqlSelectByID, id)

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

	return &post, nil
}
