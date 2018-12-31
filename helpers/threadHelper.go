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

func parentNotExists(parent int64) bool {
	if parent == 0 {
		return false
	}

	var t int
	rows := database.DB.Conn.QueryRow(`
		SELECT id
		FROM posts
		WHERE id = $1`,
		parent)

	if err := rows.Scan(&t); err != nil {
		return true
	}

	return false
}

func parentExitsInOtherThread(parent int64, threadID int) bool {
	var t int
	rows := database.DB.Conn.QueryRow(`
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

	return true
}

func AuthorExists(author string) bool {
	var nickname string
	rows := database.DB.Conn.QueryRow(`
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

	return false
}

// ThreadCreateHelper inserts thread into table THREADS
func ThreadCreateHelper(posts *models.Posts, slugOrID string) (*models.Posts, error) {
	threadByID, err := GetThreadBySlugOrID(slugOrID)
	if err != nil {
		return nil, err
	}

	if len(*posts) == 0 {
		return nil, errors.NoPostsForCreate
	}

	created := time.Now()
	insertedPosts := models.Posts{}
	for _, post := range *posts {
		var rows *pgx.Row

		if AuthorExists(post.Author) {
			return nil, errors.UserNotFound
		}

		if parentExitsInOtherThread(post.Parent, threadByID.Id) || parentNotExists(post.Parent) {
			return nil, errors.PostParentNotFound
		}

		rows = database.DB.Conn.QueryRow(`
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

	return &insertedPosts, nil
}

// GetThreadBySlugOrID selects thread by id.
func GetThreadBySlugOrID(slugOrID string) (*models.Thread, error) {
	var err error
	var thread models.Thread

	if IsNumber(slugOrID) {
		id, _ := strconv.Atoi(slugOrID)
		rows := database.DB.Conn.QueryRow(` 
			SELECT id, title, author, forum, message, votes, slug, created
			FROM threads
			WHERE id = $1`, id)

		err = rows.Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &thread.Created)
		if err != nil {
			return nil, errors.ThreadNotFound
		}
	} else {
		rows := database.DB.Conn.QueryRow(` 
			SELECT id, title, author, forum, message, votes, slug, created
			FROM threads
			WHERE slug = $1`, slugOrID)

		err = rows.Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &thread.Created)
		if err != nil {
			return nil, errors.ThreadNotFound
		}
	}

	return &thread, nil
}

// ThreadVoteHelper inserts votes into table VOTES.
func ThreadVoteHelper(v *models.Vote, slugOrID string) (*models.Thread, error) {
	_, err := UserGetOneHelper(v.Nickname)
	if err != nil {
		return nil, errors.ThreadNotFound
	}
	foundVote, _ := CheckThreadVotesByNickname(v.Nickname)
	thread, err := GetThreadBySlugOrID(slugOrID)
	if err != nil {
		return nil, err
	}

	var rows *pgx.Row
	var editedThread models.Thread
	var threadVoices int

	if foundVote == nil { // like by user did not exist
		_, err = database.DB.Conn.Exec(`
			INSERT INTO votes (nickname, voice) VALUES ($1, $2)`,
			&v.Nickname, &v.Voice)

		if err != nil {
			return nil, err
		}

		threadVoices = thread.Votes + v.Voice // counting of votes

		rows = database.DB.Conn.QueryRow(`
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

		if _, err = database.DB.Conn.Exec(`
			UPDATE votes 
			SET voice = $2
			WHERE nickname = $1`,
			&v.Nickname, &v.Voice); err != nil {
			return nil, err
		}

		threadVoices = thread.Votes + v.Voice - oldVote.Voice // recounting of votes with old voice

		rows = database.DB.Conn.QueryRow(`
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

	return &editedThread, nil
}

// ThreadGetPosts selects posts from table POSTS with filters.
func ThreadGetPosts(slugOrID string, limit, since, sort, desc []byte) (*models.Posts, error) {
	thread, err := GetThreadBySlugOrID(slugOrID)
	if err != nil {
		return nil, err
	}

	var queryRows *pgx.Rows

	if since != nil {
		if bytes.Equal([]byte("true"), desc) {
			switch string(sort) {
			case "tree":
				queryRows, err = database.DB.Conn.Query(`
					SELECT id, author, parent, message, forum, thread, created
					FROM posts
					WHERE thread = $1 AND (path < (SELECT path FROM posts WHERE id = $2::TEXT::INTEGER))
					ORDER BY path DESC
					LIMIT $3::TEXT::INTEGER`,
					thread.Id, since, limit)
			case "parent_tree":
				queryRows, err = database.DB.Conn.Query(`
					SELECT id, author, parent, message, forum, thread, created
					FROM posts
					WHERE path[1] IN (
						SELECT id
						FROM posts
						WHERE thread = $1 AND parent = 0 AND id < (SELECT path[1] FROM posts WHERE id = $2::TEXT::INTEGER)
						ORDER BY id DESC
						LIMIT $3::TEXT::INTEGER
					)
					ORDER BY path`,
					thread.Id, since, limit)
			default:
				queryRows, err = database.DB.Conn.Query(`
					SELECT id, author, parent, message, forum, thread, created
					FROM posts
					WHERE thread = $1 AND id < $2::TEXT::INTEGER
					ORDER BY id DESC
					LIMIT $3::TEXT::INTEGER`,
					thread.Id, since, limit)
			}
		} else {
			switch string(sort) {
			case "tree":
				queryRows, err = database.DB.Conn.Query(`
					SELECT id, author, parent, message, forum, thread, created
					FROM posts
					WHERE thread = $1 AND (path > (SELECT path FROM posts WHERE id = $2::TEXT::INTEGER))
					ORDER BY path
					LIMIT $3::TEXT::INTEGER`,
					thread.Id, since, limit)
			case "parent_tree":
				queryRows, err = database.DB.Conn.Query(`
					SELECT id, author, parent, message, forum, thread, created
					FROM posts
					WHERE path[1] IN (
						SELECT id
						FROM posts
						WHERE thread = $1 AND parent = 0 AND id > (SELECT path[1] FROM posts WHERE id = $2::TEXT::INTEGER)
						ORDER BY id LIMIT $3::TEXT::INTEGER
					)
					ORDER BY path`,
					thread.Id, since, limit)
			default:
				queryRows, err = database.DB.Conn.Query(`
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
				queryRows, err = database.DB.Conn.Query(`
					SELECT id, author, parent, message, forum, thread, created
					FROM posts
					WHERE thread = $1 
					ORDER BY path DESC
					LIMIT $2::TEXT::INTEGER`,
					thread.Id, limit)
			case "parent_tree":
				queryRows, err = database.DB.Conn.Query(`
					SELECT id, author, parent, message, forum, thread, created
					FROM posts
					WHERE thread = $1 AND path[1] IN (
						SELECT path[1]
						FROM posts
						WHERE thread = $1
						GROUP BY path[1]
						ORDER BY path[1] DESC
						LIMIT $2::TEXT::INTEGER
					)
					ORDER BY path[1] DESC, path;`,
					thread.Id, limit)
			default:
				queryRows, err = database.DB.Conn.Query(`
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
				queryRows, err = database.DB.Conn.Query(`
					SELECT id, author, parent, message, forum, thread, created
					FROM posts
					WHERE thread = $1 
					ORDER BY path
					LIMIT $2::TEXT::INTEGER`,
					thread.Id, limit)
			case "parent_tree":
				queryRows, err = database.DB.Conn.Query(`
					SELECT id, author, parent, message, forum, thread, created
					FROM posts
					WHERE thread = $1 AND path[1] IN (
						SELECT path[1] 
						FROM posts 
						WHERE thread = $1 
						GROUP BY path[1]
						ORDER BY path[1]
						LIMIT $2::TEXT::INTEGER
					)
					ORDER BY path`,
					thread.Id, limit)
			default:
				queryRows, err = database.DB.Conn.Query(`
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

	return &posts, nil
}

// ThreadUpdateHelper updates thread.
func ThreadUpdateHelper(thread *models.ThreadUpdate, slugOrID string) (*models.Thread, error) {
	threadFound, err := GetThreadBySlugOrID(slugOrID)
	if err != nil {
		return nil, err
	}

	updatedThread := models.Thread{}

	rows := database.DB.Conn.QueryRow(`
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
		return nil, err
	}

	return &updatedThread, nil
}
