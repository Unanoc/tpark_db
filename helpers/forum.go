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
	rows := database.DB.Conn.QueryRow(sqlInsertForum,
		&f.Slug,
		&f.Title,
		&f.User,
	)

	err := rows.Scan(&f.User)
	if err != nil {
		switch err.(pgx.PgError).Code {
		case "23505":
			forum, _ := ForumGetBySlugHelper(f.Slug)
			return forum, errors.ForumIsExist
		case "23502":
			return nil, errors.UserNotFound
		default:
			return nil, err
		}
	}

	return f, nil
}

// ForumCreateThreadHelper inserts data in table THREADS.
func ForumCreateThreadHelper(t *models.Thread) (*models.Thread, error) {
	if t.Slug != "" {
		existThread, err := GetThreadBySlugOrIDHelper(t.Slug)
		if err == nil {
			return existThread, errors.ThreadIsExist
		}
	}

	err := database.DB.Conn.QueryRow(sqlInsertThread,
		&t.Author,
		&t.Created,
		&t.Message,
		&t.Title,
		&t.Slug,
		&t.Forum,
	).Scan(
		&t.Author,
		&t.Created,
		&t.Forum,
		&t.Id,
		&t.Message,
		&t.Title,
	)

	if err != nil {
		switch err.(pgx.PgError).Code {
		case "23502":
			return nil, errors.ForumOrAuthorNotFound
		case "23503":
			return nil, errors.ForumOrAuthorNotFound
		default:
			return nil, err
		}
	}

	return t, nil
}

// ForumGetBySlugHelper selects forum by slug.
func ForumGetBySlugHelper(slug string) (*models.Forum, error) {
	forum := models.Forum{}

	err := database.DB.Conn.QueryRow(sqlSelectForumBySlug,
		slug,
	).Scan(
		&forum.Slug,
		&forum.Title,
		&forum.User,
		&forum.Posts,
		&forum.Threads,
	)

	if err != nil {
		return nil, errors.ForumNotFound
	}

	return &forum, nil
}

// ForumGetThreadsHelper selects data from THREADS with filter.
func ForumGetThreadsHelper(slug string, limit, since, desc []byte) (models.Threads, error) {
	var queryRows *pgx.Rows
	var err error

	if since != nil {
		if bytes.Equal([]byte("true"), desc) {
			queryRows, err = database.DB.Conn.Query(sqlSelectThreadsSinceDescLimit, slug, since, limit)
		} else {
			queryRows, err = database.DB.Conn.Query(sqlSelectThreadsSinceAscLimit, slug, since, limit)
		}
	} else {
		if bytes.Equal([]byte("true"), desc) {
			queryRows, err = database.DB.Conn.Query(sqlSelectThreadDescLimit, slug, limit)
		} else {
			queryRows, err = database.DB.Conn.Query(sqlSelectThreadAscLimit, slug, limit)
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
			&thread.Votes,
		)

		threads = append(threads, &thread)
	}

	if len(threads) == 0 {
		_, err := ForumGetBySlugHelper(slug)
		if err != nil {
			return nil, errors.ForumNotFound
		}
	}

	return threads, nil
}

// ForumGetUsersHelper selects users of forum from table USERS.
func ForumGetUsersHelper(slug string, limit, since, desc []byte) (*models.Users, error) {
	_, err := ForumGetBySlugHelper(slug)
	if err != nil {
		return nil, err
	}

	var queryRows *pgx.Rows

	if since != nil {
		if bytes.Equal([]byte("true"), desc) {
			queryRows, err = database.DB.Conn.Query(sqlSelectUsersDescSinceLimit, slug, since, limit)
		} else {
			queryRows, err = database.DB.Conn.Query(sqlSelectUsersAscSinceLimit, slug, since, limit)
		}
	} else {
		if bytes.Equal([]byte("true"), desc) {
			queryRows, err = database.DB.Conn.Query(sqlSelectUsersDescLimit, slug, limit)
		} else {
			queryRows, err = database.DB.Conn.Query(sqlSelectUsersAscLimit, slug, limit)
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
		_, err := ForumGetBySlugHelper(slug)
		if err != nil {
			return nil, errors.UserNotFound
		}
	}

	return &users, nil
}
