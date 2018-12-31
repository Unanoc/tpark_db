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
