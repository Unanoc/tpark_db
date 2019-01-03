package helpers

import (
	"tpark_db/database"
	"tpark_db/errors"
	"tpark_db/models"

	"github.com/jackc/pgx"
)

// UserCreateHelper inserts user into table USERS.
func UserCreateHelper(u *models.User) (models.Users, error) {
	rows, err := database.DB.Conn.Exec(sqlInsertUser, &u.Nickname, &u.Fullname, &u.About, &u.Email)

	if err != nil {
		return nil, err
	}

	if rows.RowsAffected() == 0 { // if it returns 0 - user existed, else user was created
		users := models.Users{}
		queryRows, err := database.DB.Conn.Query(sqlSelectUserByEmailOrNickName, &u.Email, &u.Nickname)
		defer queryRows.Close()

		if err != nil {
			return nil, err
		}

		for queryRows.Next() {
			user := models.User{}
			queryRows.Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email)
			users = append(users, &user)
		}
		return users, errors.UserIsExist
	}

	return nil, nil
}

// UserUpdateHelper updates user.
func UserUpdateHelper(user *models.User) error {
	err := database.DB.Conn.QueryRow(sqlUpdateUser,
		&user.Nickname,
		&user.Fullname,
		&user.About,
		&user.Email,
	).Scan(
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

// UserGetOneHelper selects user by username from table USERS.
func UserGetOneHelper(nickname string) (*models.User, error) {
	user := models.User{}

	err := database.DB.Conn.QueryRow(sqlSelectUserByNickname, nickname).Scan(
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
