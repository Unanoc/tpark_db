package helpers

import (
	"tpark_db/database"
	"tpark_db/errors"
	"tpark_db/models"

	"github.com/jackc/pgx"
)

// UserCreateHelper inserts user into table USERS.
func UserCreateHelper(u *models.User) (models.Users, error) {
	rows, err := database.DB.Conn.Exec(`
		INSERT
		INTO users ("nickname", "fullname", "about", "email")
		VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`,
		&u.Nickname, &u.Fullname, &u.About, &u.Email)

	if err != nil {
		return nil, err
	}

	if rows.RowsAffected() == 0 { // if it returns 0 - user existed, else user was created
		users := models.Users{}
		queryRows, err := database.DB.Conn.Query(`
			SELECT "nickname", "fullname", "about", "email"
			FROM users
			WHERE "email" = $1 OR "nickname" = $2`,
			&u.Email, &u.Nickname)

		defer queryRows.Close()

		if err != nil {
			return nil, err
		}

		for queryRows.Next() {
			user := models.User{}
			queryRows.Scan(&user.Nickname, &user.Fullname,
				&user.About, &user.Email)
			users = append(users, &user)
		}
		return users, errors.UserIsExist
	}

	return nil, nil
}

// UserGetOneHelper selects user by username from table USERS.
func UserGetOneHelper(username string) (*models.User, error) {
	user := models.User{}

	err := database.DB.Conn.QueryRow(`
		SELECT "nickname", "fullname", "about", "email"
		FROM users
		WHERE "nickname" = $1`,
		username).Scan(
		&user.Nickname,
		&user.Fullname,
		&user.About,
		&user.Email,
	)

	if err != nil {
		return nil, errors.UserNotFound
	}

	return &user, nil
}

// UserUpdateHelper updates user.
func UserUpdateHelper(user *models.User) error {
	rows := database.DB.Conn.QueryRow(`
		UPDATE users
		SET fullname = coalesce(nullif($2, ''), fullname),
			about    = coalesce(nullif($3, ''), about),
			email    = coalesce(nullif($4, ''), email)
		WHERE "nickname" = $1
		RETURNING fullname, about, email, nickname`,
		&user.Nickname, &user.Fullname, &user.About, &user.Email)

	err := rows.Scan(
		&user.Fullname,
		&user.About,
		&user.Email,
		&user.Nickname,
	)

	if err != nil {
		if _, ok := err.(pgx.PgError); ok {
			return errors.UserUpdateConflict
		}
		return errors.UserNotFound
	}

	return nil
}
