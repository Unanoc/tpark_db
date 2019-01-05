package helpers

import (
	"tpark_db/database"
	"tpark_db/errors"
	"tpark_db/models"
	"unicode"

	"github.com/jackc/pgx"
)

// IsNumber checks if string is number.
func IsNumber(s string) bool {
	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}

// CheckThreadVotesByNickname chechs if vote exists.
func CheckThreadVotesByNickname(nickname string) (*models.Vote, error) {
	rows := database.DB.Conn.QueryRow(` 
		SELECT nickname, voice
		FROM votes
		WHERE nickname = $1`,
		nickname)

	var vote models.Vote
	err := rows.Scan(&vote.Nickname, &vote.Voice)
	if err != nil {
		return nil, err
	}

	return &vote, nil
}

func parentNotExists(parent int64) error {
	if parent == 0 {
		return nil
	}

	var t int
	rows := database.DB.Conn.QueryRow(`
		SELECT id
		FROM posts
		WHERE id = $1`,
		parent)

	if err := rows.Scan(&t); err != nil {
		if err == pgx.ErrNoRows {
			return errors.PostParentNotFound
		}
	}

	return nil
}

func parentExitsInOtherThread(parent int64, threadID int) error {
	var t int
	rows := database.DB.Conn.QueryRow(`
		SELECT id
		FROM posts
		WHERE id = $1 AND thread IN (SELECT id FROM threads WHERE thread <> $2)`,
		parent, threadID)

	if err := rows.Scan(&t); err != nil {
		if err == pgx.ErrNoRows {
			return nil
		}
	}

	return errors.PostParentNotFound
}

func authorExists(nickname string) error {
	var user models.User
	rows := database.DB.Conn.QueryRow(sqlSelectUserByNickname, nickname)

	if err := rows.Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email); err != nil {
		if err == pgx.ErrNoRows {
			return errors.UserNotFound
		}
	}

	return nil
}
