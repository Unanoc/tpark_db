package helpers

import (
	"strconv"
	"time"
	"tpark_db/database"
	"tpark_db/errors"
	"tpark_db/models"

	"github.com/jackc/pgx"
)

func ThreadCreateHelper(posts *models.Posts, slugOrID string) (*models.Posts, error) {
	if len(*posts) == 0 {
		return nil, errors.NoPostsForCreate
	}

	threadByID, err := GetThreadBySlugOrId(slugOrID)
	if err != nil {
		return nil, err
	}

	tx := database.StartTransaction()
	defer tx.Rollback()

	created := time.Now()
	insertedPosts := models.Posts{}
	for _, post := range *posts {
		var rows *pgx.Row

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

func ThreadGetPosts(slugOrID string, limit, since, sort, desc []byte) (*models.Threads, error) {
	// thread, err := GetThreadBySlugOrId(slugOrID)
	// if err != nil {
	// 	return nil, err
	// }

	// tx := database.StartTransaction()
	// defer tx.Rollback()
	// var queryRows *pgx.Rows
	// var err error

	// if since != nil {
	// 	if bytes.Equal([]byte("true"), desc) {
	// 		queryRows, err = tx.Query(`
	// 			SELECT nickname, fullname, about, email
	// 			FROM users
	// 			WHERE forum = $1 AND created <= $2::TEXT::TIMESTAMPTZ
	// 			ORDER BY created DESC
	// 			LIMIT $3::TEXT::INTEGER`,
	// 			slug, since, limit)
	// 	} else {
	// 		queryRows, err = tx.Query(`
	// 			SELECT nickname, fullname, about, email
	// 			FROM users
	// 			WHERE forum = $1 AND created >= $2::TEXT::TIMESTAMPTZ
	// 			ORDER BY created
	// 			LIMIT $3::TEXT::INTEGER`,
	// 			slug, since, limit)
	// 	}
	// } else {
	// 	if bytes.Equal([]byte("true"), desc) {
	// 		queryRows, err = tx.Query(`
	// 			SELECT nickname, fullname, about, email
	// 			FROM users
	// 			WHERE forum = $1
	// 			ORDER BY created DESC
	// 			LIMIT $2::TEXT::INTEGER`,
	// 			slug, limit)
	// 	} else {
	// 		queryRows, err = tx.Query(`
	// 			SELECT author, created, forum, id, message, slug, title, votes
	// 			FROM threads
	// 			WHERE forum = $1
	// 			ORDER BY created
	// 			LIMIT $2::TEXT::INTEGER`,
	// 			slug, limit)
	// 	}
	// }
	// defer queryRows.Close()

	// if err != nil {
	// 	return nil, errors.UserNotFound
	// }

	// users := models.Users{}
	// for queryRows.Next() {
	// 	user := models.User{}

	// 	if err = queryRows.Scan(&user.Nickname, &user.Fullname, &user.About,
	// 		&user.Email); err != nil {
	// 		fmt.Println(err)
	// 	}
	// 	users = append(users, &user)
	// }

	// if len(users) == 0 {
	// 	_, err := ForumGetBySlug(slug)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		return nil, errors.UserNotFound
	// 	}
	// }

	// database.CommitTransaction(tx)
	return nil, errors.ThreadNotFound
}
