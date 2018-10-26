package helpers

import (
	"tpark_db/database"
	"tpark_db/models"
	"unicode"
)

func IsNumber(s string) bool {
	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}

func CheckThreadVotesByNickname(nickname string) (*models.Vote, error) {
	tx := database.StartTransaction()
	defer tx.Rollback()

	rows := tx.QueryRow(` 
		SELECT nickname, voice
		FROM votes
		WHERE nickname = $1`,
		nickname)

	var vote models.Vote
	err := rows.Scan(&vote.Nickname, &vote.Voice)
	if err != nil {
		return nil, err
	}

	database.CommitTransaction(tx)
	return &vote, nil
}
