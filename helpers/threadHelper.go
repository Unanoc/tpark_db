package helpers

import (
	"bytes"
	"fmt"
	"strconv"
	"time"
	"tpark_db/database"
	"tpark_db/errors"
	"tpark_db/models"

	"github.com/jackc/pgx"
)

func ParentNotExists(parent int64) bool {
	tx := database.StartTransaction()
	defer tx.Rollback()

	if parent == 0 {
		return false
	}

	var t int
	rows := tx.QueryRow(`
		SELECT id
		FROM posts
		WHERE id = $1`,
		parent)

	if err := rows.Scan(&t); err != nil {
		return true
	}
	database.CommitTransaction(tx)
	return false
}

func ParentExitsInOtherThread(parent int64, threadID int) bool {
	tx := database.StartTransaction()
	defer tx.Rollback()

	var t int
	rows := tx.QueryRow(`
		SELECT id
		FROM posts
		WHERE id = $1 AND thread IN (SELECT id FROM threads WHERE thread <> $2)`,
		parent, threadID)

	if err := rows.Scan(&t); err != nil {
		if err.Error() == "no rows in result set" {
			return false
		}
		return true
	}

	database.CommitTransaction(tx)
	return true
}

func AuthorExists(author string) bool {
	tx := database.StartTransaction()
	defer tx.Rollback()

	var nickname string
	rows := tx.QueryRow(`
		SELECT nickname
		FROM users
		WHERE nickname = $1`,
		author)

	if err := rows.Scan(&nickname); err != nil {
		if err.Error() == "no rows in result set" {
			return true
		}
		return false
	}

	database.CommitTransaction(tx)
	return false
}

func ThreadCreateHelper(posts *models.Posts, slugOrID string) (*models.Posts, error) {
	threadByID, err := GetThreadBySlugOrId(slugOrID)
	if err != nil {
		return nil, err
	}

	if len(*posts) == 0 {
		return nil, errors.NoPostsForCreate
	}

	tx := database.StartTransaction()
	defer tx.Rollback()

	created := time.Now()
	insertedPosts := models.Posts{}
	for _, post := range *posts {
		var rows *pgx.Row

		if AuthorExists(post.Author) {
			return nil, errors.UserNotFound
		}

		if ParentExitsInOtherThread(post.Parent, threadByID.Id) || ParentNotExists(post.Parent) {
			return nil, errors.PostParentNotFound
		}

		rows = tx.QueryRow(`
				INSERT
				INTO posts (author, created, message, thread, parent, forum, path)
				VALUES ($1, $2, $3, $4, $5, $6, (SELECT path FROM posts WHERE id = $5) || (select currval(pg_get_serial_sequence('posts', 'id'))) )
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

	database.CommitTransaction(tx)
	return &insertedPosts, nil
}

func GetThreadBySlugOrId(slugOrId string) (*models.Thread, error) {
	var err error
	var thread models.Thread

	tx := database.StartTransaction()
	defer tx.Rollback()

	if IsNumber(slugOrId) {
		id, _ := strconv.Atoi(slugOrId)
		rows := tx.QueryRow(` 
			SELECT id, title, author, forum, message, votes, slug, created
			FROM threads
			WHERE id = $1`, id)

		err = rows.Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &thread.Created)
		if err != nil {
			return nil, errors.ThreadNotFound
		}
	} else {
		rows := tx.QueryRow(` 
			SELECT id, title, author, forum, message, votes, slug, created
			FROM threads
			WHERE slug = $1`, slugOrId)

		err = rows.Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &thread.Created)
		if err != nil {
			return nil, errors.ThreadNotFound
		}
	}

	database.CommitTransaction(tx)
	return &thread, nil
}

func ThreadVoteHelper(v *models.Vote, slugOrID string) (*models.Thread, error) {
	tx := database.StartTransaction()
	defer tx.Rollback()

	_, err := UserGetOneHelper(v.Nickname)
	if err != nil {
		return nil, errors.ThreadNotFound
	}
	foundVote, _ := CheckThreadVotesByNickname(v.Nickname)
	thread, err := GetThreadBySlugOrId(slugOrID)
	if err != nil {
		return nil, err
	}

	var rows *pgx.Row
	var editedThread models.Thread
	var threadVoices int

	if foundVote == nil { // like by user did not exist
		_, err = tx.Exec(`
			INSERT INTO votes (nickname, voice) VALUES ($1, $2)`,
			&v.Nickname, &v.Voice)

		if err != nil {
			return nil, err
		}

		threadVoices = thread.Votes + v.Voice // counting of votes

		rows = tx.QueryRow(`
			UPDATE threads
			SET votes = $1
			WHERE slug = $2
			RETURNING id, title, author, forum, message, votes, slug, created`, &threadVoices, &thread.Slug,
		)
		err = rows.Scan(
			&editedThread.Id,
			&editedThread.Title,
			&editedThread.Author,
			&editedThread.Forum,
			&editedThread.Message,
			&editedThread.Votes,
			&editedThread.Slug,
			&editedThread.Created,
		)

		if err != nil {
			return nil, err
		}
	} else {
		oldVote, _ := CheckThreadVotesByNickname(v.Nickname)

		if _, err = tx.Exec(`
			UPDATE votes 
			SET voice = $2
			WHERE nickname = $1`,
			&v.Nickname, &v.Voice); err != nil {
			return nil, err
		}

		threadVoices = thread.Votes + v.Voice - oldVote.Voice // recounting of votes with old voice

		rows = tx.QueryRow(`
			UPDATE threads
			SET votes = $1
			WHERE slug = $2
			RETURNING id, title, author, forum, message, votes, slug, created`, &threadVoices, &thread.Slug,
		)

		err = rows.Scan(
			&editedThread.Id,
			&editedThread.Title,
			&editedThread.Author,
			&editedThread.Forum,
			&editedThread.Message,
			&editedThread.Votes,
			&editedThread.Slug,
			&editedThread.Created,
		)

		if err != nil {
			return nil, err
		}
	}

	database.CommitTransaction(tx)
	return &editedThread, nil
}

func ThreadGetPosts(slugOrID string, limit, since, sort, desc []byte) (*models.Posts, error) {
	thread, err := GetThreadBySlugOrId(slugOrID)
	if err != nil {
		return nil, err
	}

	tx := database.StartTransaction()
	defer tx.Rollback()
	var queryRows *pgx.Rows

	if since != nil {
		if bytes.Equal([]byte("true"), desc) {
			switch string(sort) {
			// case "tree":
			// 	//TODO
			// case "parent_tree":
			// 	//TODO
			default:
				queryRows, err = tx.Query(`
					SELECT id, author, parent, message, forum, thread, created
					FROM posts
					WHERE thread = $1 AND id < $2
					ORDER BY id DESC
					LIMIT $3::TEXT::INTEGER`,
					thread.Id, since, limit)
			}
		} else {
			switch string(sort) {
			case "tree":
				queryRows, err = tx.Query(`
					SELECT id, author, parent, message, forum, thread, created
					FROM posts
					WHERE thread = $1 AND (path > (SELECT path FROM posts WHERE id = $2::TEXT::INTEGER))
					ORDER BY path
					LIMIT $3::TEXT::INTEGER`,
					thread.Id, since, limit)
			case "parent_tree":
				queryRows, err = tx.Query(`
				SELECT id, author, parent, message, forum, thread, created
				FROM posts
				WHERE thread = $1 AND parent IN (
					SELECT parent 
					FROM posts 
					WHERE thread = $1 id > $2::TEXT::INTEGER
					ORDER BY parent 
					LIMIT $3::TEXT::INTEGER
				)
				ORDER BY path`,
					thread.Id, since, limit)
			default:
				queryRows, err = tx.Query(`
					SELECT id, author, parent, message, forum, thread, created
					FROM posts
					WHERE thread = $1 AND id > $2::TEXT::INTEGER
					ORDER BY id
					LIMIT $3::TEXT::INTEGER`,
					thread.Id, since, limit)
			}
		}
	} else {
		if bytes.Equal([]byte("true"), desc) {
			switch string(sort) {
			case "tree":
				queryRows, err = tx.Query(`
					SELECT id, author, parent, message, forum, thread, created
					FROM posts
					WHERE thread = $1 
					ORDER BY path DESC
					LIMIT $2::TEXT::INTEGER`,
					thread.Id, limit)
			case "parent_tree":
				queryRows, err = tx.Query(`
					SELECT id, author, parent, message, forum, thread, created
					FROM posts 
					WHERE thread = $1 
					ORDER BY path[1] DESC, array_length(path, 1), path[2]
					LIMIT $2::TEXT::INTEGER`,
					thread.Id, limit)
			default:
				queryRows, err = tx.Query(`
					SELECT id, author, parent, message, forum, thread, created
					FROM posts
					WHERE thread = $1
					ORDER BY id DESC
					LIMIT $2::TEXT::INTEGER`,
					thread.Id, limit)
			}
		} else {
			switch string(sort) {
			case "tree":
				queryRows, err = tx.Query(`
					SELECT id, author, parent, message, forum, thread, created
					FROM posts
					WHERE thread = $1 
					ORDER BY path
					LIMIT $2::TEXT::INTEGER`,
					thread.Id, limit)
			case "parent_tree":
				queryRows, err = tx.Query(`
					SELECT id, author, parent, message, forum, thread, created
					FROM posts
					WHERE thread = $1 AND parent IN (
						SELECT parent 
						FROM posts 
						WHERE thread = $1 
						ORDER BY path 
						LIMIT $2::TEXT::INTEGER
					)
					ORDER BY path`,
					thread.Id, limit)
			default:
				queryRows, err = tx.Query(`
					SELECT id, author, parent, message, forum, thread, created
					FROM posts
					WHERE thread = $1 
					ORDER BY id
					LIMIT $2::TEXT::INTEGER`,
					thread.Id, limit)
			}
		}
	}
	defer queryRows.Close()

	if err != nil {
		return nil, err
	}

	posts := models.Posts{}
	for queryRows.Next() {
		post := models.Post{}

		if err = queryRows.Scan(
			&post.Id,
			&post.Author,
			&post.Parent,
			&post.Message,
			&post.Forum,
			&post.Thread,
			&post.Created,
		); err != nil {
			fmt.Println(err)
		}
		posts = append(posts, &post)
	}

	database.CommitTransaction(tx)
	return &posts, nil
}

func ThreadUpdateHelper(thread *models.ThreadUpdate, slugOrID string) (*models.Thread, error) {
	threadFound, err := GetThreadBySlugOrId(slugOrID)
	if err != nil {
		return nil, err
	}

	tx := database.StartTransaction()
	defer tx.Rollback()

	updatedThread := models.Thread{}

	rows := tx.QueryRow(`
		UPDATE threads
		SET title = coalesce(nullif($2, ''), title),
			message = coalesce(nullif($3, ''), message)
		WHERE slug = $1
		RETURNING id, title, author, forum, message, votes, slug, created`,
		&threadFound.Slug,
		&thread.Title,
		&thread.Message)

	err = rows.Scan(
		&updatedThread.Id,
		&updatedThread.Title,
		&updatedThread.Author,
		&updatedThread.Forum,
		&updatedThread.Message,
		&updatedThread.Votes,
		&updatedThread.Slug,
		&updatedThread.Created,
	)

	if err != nil {
		// if _, ok := err.(pgx.PgError); ok {
		// 	return nil, errors.UserUpdateConflict
		// }
		// return nil, errors.UserNotFound
		return nil, err
	}

	database.CommitTransaction(tx)
	return &updatedThread, nil
}
