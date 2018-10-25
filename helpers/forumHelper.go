package helpers

import (
	"bytes"
	"fmt"
	"tpark_db/database"
	"tpark_db/errors"
	"tpark_db/models"

	"github.com/jackc/pgx"
)

func ForumCreateHelper(f *models.Forum) (*models.Forum, error) {
	tx := database.StartTransaction()
	defer tx.Rollback()

	rows := tx.QueryRow(`
		INSERT
		INTO forums (slug, title, "user")
		VALUES ($1, $2, (SELECT nickname FROM users WHERE nickname = $3)) 
		RETURNING "user"`,
		&f.Slug,
		&f.Title,
		&f.User)

	err := rows.Scan(&f.User)

	if err != nil {
		switch err.(pgx.PgError).Code {
		case "23505":
			forum, _ := ForumGetBySlug(f.Slug)
			return forum, errors.ForumIsExist
		case "23502":
			return nil, errors.UserNotFound
		default:
			return nil, err
		}
	}

	database.CommitTransaction(tx)
	return f, nil
}

func ForumGetBySlug(slug string) (*models.Forum, error) {
	tx := database.StartTransaction()
	defer tx.Rollback()

	forum := models.Forum{}

	err := tx.QueryRow(`
		SELECT slug, title, "user", posts, threads
		FROM forums
		WHERE slug = $1`,
		slug).Scan(
		&forum.Slug,
		&forum.Title,
		&forum.User,
		&forum.Posts,
		&forum.Threads)

	if err != nil {
		return &forum, errors.ForumNotFound
	}

	database.CommitTransaction(tx)
	return &forum, nil
}

func ForumCreateThreadHelper(t *models.Thread) (*models.Thread, error) {
	if t.Slug != "" {
		existThread, err := GetThreadBySlug(t.Slug)
		if err == nil {
			return existThread, errors.ThreadIsExist
		}
	}

	tx := database.StartTransaction()
	defer tx.Rollback()

	rows := tx.QueryRow(`
		INSERT
		INTO threads (author, created, message, title, slug, forum)
		VALUES ($1, $2, $3, $4, $5, (SELECT slug FROM forums WHERE slug = $6)) 
		RETURNING author, created, forum, id, message, title`,
		&t.Author,
		&t.Created,
		&t.Message,
		&t.Title,
		&t.Slug,
		&t.Forum)

	err := rows.Scan(&t.Author, &t.Created, &t.Forum, &t.Id, &t.Message, &t.Title)
	if err != nil {
		switch err.(pgx.PgError).Code {
		case "23503":
			return nil, errors.ForumOrAuthorNotFound
		case "23502":
			return nil, errors.ForumOrAuthorNotFound
		default:
			return nil, err
		}
	}

	database.CommitTransaction(tx)
	return t, nil
}

func ForumGetThreadsHelper(slug string, limit, since, desc []byte) (models.Threads, error) {
	tx := database.StartTransaction()
	defer tx.Rollback()
	var queryRows *pgx.Rows
	var err error

	if since != nil {
		if bytes.Equal([]byte("true"), desc) {
			queryRows, err = tx.Query(`
				SELECT author, created, forum, id, message, slug, title, votes
				FROM threads
				WHERE forum = $1 AND created <= $2::TEXT::TIMESTAMPTZ
				ORDER BY created DESC
				LIMIT $3::TEXT::INTEGER`,
				slug, since, limit)
		} else {
			queryRows, err = tx.Query(`
				SELECT author, created, forum, id, message, slug, title, votes
				FROM threads
				WHERE forum = $1 AND created >= $2::TEXT::TIMESTAMPTZ
				ORDER BY created
				LIMIT $3::TEXT::INTEGER`,
				slug, since, limit)
		}
	} else {
		if bytes.Equal([]byte("true"), desc) {
			queryRows, err = tx.Query(`
				SELECT author, created, forum, id, message, slug, title, votes
				FROM threads
				WHERE forum = $1
				ORDER BY created DESC
				LIMIT $2::TEXT::INTEGER`,
				slug, limit)
		} else {
			queryRows, err = tx.Query(`
				SELECT author, created, forum, id, message, slug, title, votes
				FROM threads
				WHERE forum = $1
				ORDER BY created
				LIMIT $2::TEXT::INTEGER`,
				slug, limit)
		}
	}
	defer queryRows.Close()

	if err != nil {
		return nil, errors.ForumNotFound
	}

	threads := models.Threads{}
	for queryRows.Next() {
		thread := models.Thread{}

		if err = queryRows.Scan(&thread.Author, &thread.Created, &thread.Forum,
			&thread.Id, &thread.Message, &thread.Slug, &thread.Title,
			&thread.Votes); err != nil {
			fmt.Println(err)
		}
		threads = append(threads, &thread)
	}

	if len(threads) == 0 {
		_, err := ForumGetBySlug(slug)
		if err != nil {
			fmt.Println(err)
			return nil, errors.ForumNotFound
		}
	}

	database.CommitTransaction(tx)
	return threads, nil
}

func ForumGetUsersHelper(slug string, limit, since, desc []byte) (models.Users, error) {
	tx := database.StartTransaction()
	defer tx.Rollback()
	var queryRows *pgx.Rows
	var err error

	if since != nil {
		if bytes.Equal([]byte("true"), desc) {
			queryRows, err = tx.Query(`
				SELECT nickname, fullname, about, email
				FROM users
				WHERE forum = $1 AND created <= $2::TEXT::TIMESTAMPTZ
				ORDER BY created DESC
				LIMIT $3::TEXT::INTEGER`,
				slug, since, limit)
		} else {
			queryRows, err = tx.Query(`
				SELECT nickname, fullname, about, email
				FROM users
				WHERE forum = $1 AND created >= $2::TEXT::TIMESTAMPTZ
				ORDER BY created
				LIMIT $3::TEXT::INTEGER`,
				slug, since, limit)
		}
	} else {
		if bytes.Equal([]byte("true"), desc) {
			queryRows, err = tx.Query(`
				SELECT nickname, fullname, about, email
				FROM users
				WHERE forum = $1
				ORDER BY created DESC
				LIMIT $2::TEXT::INTEGER`,
				slug, limit)
		} else {
			queryRows, err = tx.Query(`
				SELECT author, created, forum, id, message, slug, title, votes
				FROM threads
				WHERE forum = $1
				ORDER BY created
				LIMIT $2::TEXT::INTEGER`,
				slug, limit)
		}
	}
	defer queryRows.Close()

	if err != nil {
		return nil, errors.UserNotFound
	}

	users := models.Users{}
	for queryRows.Next() {
		user := models.User{}

		if err = queryRows.Scan(&user.Nickname, &user.Fullname, &user.About,
			&user.Email); err != nil {
			fmt.Println(err)
		}
		users = append(users, &user)
	}

	if len(users) == 0 {
		_, err := ForumGetBySlug(slug)
		if err != nil {
			fmt.Println(err)
			return nil, errors.UserNotFound
		}
	}

	database.CommitTransaction(tx)
	return users, nil
}
