package helpers

import (
	"tpark_db/database"
	"tpark_db/models"
	"unicode"
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

func authorExists(nickname string) bool {
	var user models.User
	rows := database.DB.Conn.QueryRow(sqlSelectUserByNickname, nickname)

	if err := rows.Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email); err != nil {
		if err.Error() == "no rows in result set" {
			return true
		}
		return false
	}

	return false
}
