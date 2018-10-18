package helpers

import (
	"log"
	"tpark_db/database"
	"tpark_db/errors"
	"tpark_db/models"
)

func UserCreateHelper(u *models.User) (models.Users, error) {
	tx := database.StartTransaction()
	defer tx.Rollback()

	rows, err := tx.Exec(`
		INSERT
		INTO users ("nickname", "fullname", "about", "email")
		VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`,
		&u.Nickname, &u.Fullname, &u.About, &u.Email)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	if rows.RowsAffected() == 0 { // if it returns 0 - user existed, else user was created
		users := models.Users{}
		queryRows, err := tx.Query(`
			SELECT "nickname", "fullname", "about", "email"
			FROM users
			WHERE "email" = $1 OR "nickname" = $2`,
			&u.Email, &u.Nickname)

		defer queryRows.Close()

		if err != nil {
			log.Println(err)
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

	database.CommitTransaction(tx)
	return nil, nil
}

func UserGetOneHelper(username string) (models.User, error) {
	tx := database.StartTransaction()
	defer tx.Rollback()

	user := models.User{}

	if err := tx.QueryRow(`
		SELECT "nickname", "fullname", "about", "email"
		FROM users
		WHERE "nickname" = $1`,
		username).Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email); err != nil {
		log.Println(err)
		return user, errors.UserNotFound
	}

	database.CommitTransaction(tx)
	return user, nil
}

func UserUpdateHelper(user *models.User) (*models.User, error) {
	tx := database.StartTransaction()
	defer tx.Rollback()
	// TODO
	database.CommitTransaction(tx)
	return user, nil
}
