package helpers

import (
	"strconv"
	"tpark_db/database"
	"tpark_db/errors"
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

func GetThreadBySlugOrId(slugOrId string) (*models.Thread, error) {
	var err error
	var thread models.Thread

	tx := database.StartTransaction()
	defer tx.Rollback()

	if IsNumber(slugOrId) {
		id, _ := strconv.Atoi(slugOrId)
		rows := tx.QueryRow(` 
			SELECT id, title, author, forum, message, votes, slug, created
			FROM threads
			WHERE id = $1`, id)

		err = rows.Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &thread.Created)
		if err != nil {
			return nil, errors.ThreadNotFound
		}
	} else {
		rows := tx.QueryRow(` 
			SELECT id, title, author, forum, message, votes, slug, created
			FROM threads
			WHERE slug = $1`, slugOrId)

		err = rows.Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &thread.Created)
		if err != nil {
			return nil, errors.ThreadNotFound
		}
	}

	database.CommitTransaction(tx)
	return &thread, nil
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
