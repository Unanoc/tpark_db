package helpers

import (
	"bytes"
	"tpark_db/database"
	"tpark_db/errors"
	"tpark_db/models"

	"github.com/jackc/pgx"
)

// ForumCreateHelper inserts data in table FORUMS.
func ForumCreateHelper(f *models.Forum) (*models.Forum, error) {
	rows := database.DB.Conn.QueryRow(`
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

	return f, nil
}

// ForumGetBySlug selects forum by slug.
func ForumGetBySlug(slug string) (*models.Forum, error) {
	forum := models.Forum{}

	err := database.DB.Conn.QueryRow(`
		SELECT slug, title, "user", 
			(SELECT COUNT(*) FROM posts WHERE forum = $1), 
			(SELECT COUNT(*) FROM threads WHERE forum = $1)
		FROM forums
		WHERE slug = $1`,
		slug).Scan(
		&forum.Slug,
		&forum.Title,
		&forum.User,
		&forum.Posts,
		&forum.Threads)

	if err != nil {
		return nil, errors.ForumNotFound
	}

	return &forum, nil
}

// ForumCreateThreadHelper inserts data in table THREADS.
func ForumCreateThreadHelper(t *models.Thread) (*models.Thread, error) {
	if t.Slug != "" {
		existThread, err := GetThreadBySlugOrID(t.Slug)
		if err == nil {
			return existThread, errors.ThreadIsExist
		}
	}

	rows := database.DB.Conn.QueryRow(`
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

	return t, nil
}

// ForumGetThreadsHelper selects data from THREADS with filter.
func ForumGetThreadsHelper(slug string, limit, since, desc []byte) (models.Threads, error) {
	var queryRows *pgx.Rows
	var err error

	if since != nil {
		if bytes.Equal([]byte("true"), desc) {
			queryRows, err = database.DB.Conn.Query(`
				SELECT author, created, forum, id, message, slug, title, votes
				FROM threads
				WHERE forum = $1 AND created <= $2::TEXT::TIMESTAMPTZ
				ORDER BY created DESC
				LIMIT $3::TEXT::INTEGER`,
				slug, since, limit)
		} else {
			queryRows, err = database.DB.Conn.Query(`
				SELECT author, created, forum, id, message, slug, title, votes
				FROM threads
				WHERE forum = $1 AND created >= $2::TEXT::TIMESTAMPTZ
				ORDER BY created
				LIMIT $3::TEXT::INTEGER`,
				slug, since, limit)
		}
	} else {
		if bytes.Equal([]byte("true"), desc) {
			queryRows, err = database.DB.Conn.Query(`
				SELECT author, created, forum, id, message, slug, title, votes
				FROM threads
				WHERE forum = $1
				ORDER BY created DESC
				LIMIT $2::TEXT::INTEGER`,
				slug, limit)
		} else {
			queryRows, err = database.DB.Conn.Query(`
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

		_ = queryRows.Scan(
			&thread.Author,
			&thread.Created,
			&thread.Forum,
			&thread.Id,
			&thread.Message,
			&thread.Slug,
			&thread.Title,
			&thread.Votes)

		threads = append(threads, &thread)
	}

	if len(threads) == 0 {
		_, err := ForumGetBySlug(slug)
		if err != nil {
			return nil, errors.ForumNotFound
		}
	}

	return threads, nil
}

// ForumGetUsersHelper selects users of forum from table USERS.
func ForumGetUsersHelper(slug string, limit, since, desc []byte) (*models.Users, error) {
	_, err := ForumGetBySlug(slug)
	if err != nil {
		return nil, err
	}

	var queryRows *pgx.Rows

	if since != nil {
		if bytes.Equal([]byte("true"), desc) {
			queryRows, err = database.DB.Conn.Query(`
				SELECT nickname, fullname, about, email
				FROM users
				WHERE nickname IN (
						SELECT author FROM threads WHERE forum = $1
						UNION
						SELECT author FROM posts WHERE forum = $1
					) 
					AND LOWER(nickname) < LOWER($2::TEXT)
				ORDER BY nickname DESC
				LIMIT $3::TEXT::INTEGER`,
				slug, since, limit)
		} else {
			queryRows, err = database.DB.Conn.Query(`
				SELECT nickname, fullname, about, email
				FROM users
				WHERE nickname IN (
						SELECT author FROM threads WHERE forum = $1
						UNION
						SELECT author FROM posts WHERE forum = $1
					)  
					AND LOWER(nickname) > LOWER($2::TEXT)
				ORDER BY nickname
				LIMIT $3::TEXT::INTEGER`,
				slug, since, limit)
		}
	} else {
		if bytes.Equal([]byte("true"), desc) {
			queryRows, err = database.DB.Conn.Query(`
				SELECT nickname, fullname, about, email
				FROM users
				WHERE nickname IN (
						SELECT author FROM threads WHERE forum = $1
						UNION
						SELECT author FROM posts WHERE forum = $1
					) 
				ORDER BY nickname DESC
				LIMIT $2::TEXT::INTEGER`,
				slug, limit)
		} else {
			queryRows, err = database.DB.Conn.Query(`
				SELECT nickname, fullname, about, email
				FROM users
				WHERE nickname IN (
						SELECT author FROM threads WHERE forum = $1
						UNION
						SELECT author FROM posts WHERE forum = $1
					) 
				ORDER BY nickname
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

		err = queryRows.Scan(
			&user.Nickname,
			&user.Fullname,
			&user.About,
			&user.Email)

		users = append(users, &user)
	}

	if len(users) == 0 {
		_, err := ForumGetBySlug(slug)
		if err != nil {
			return nil, errors.UserNotFound
		}
	}

	return &users, nil
}
