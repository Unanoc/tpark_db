package helpers

import (
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

func GetThreadBySlug(slug string) (*models.Thread, error) {
	tx := database.StartTransaction()
	defer tx.Rollback()

	rows := tx.QueryRow(` 
		SELECT id, title, author, forum, message, votes, slug, created
		FROM threads
		WHERE slug = $1
	`, slug)

	var thread models.Thread
	err := rows.Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &thread.Created)
	if err != nil {
		return nil, errors.ThreadNotFound
	}

	database.CommitTransaction(tx)
	return &thread, nil
}

func GetThreadById(id int) (*models.Thread, error) {
	tx := database.StartTransaction()
	defer tx.Rollback()

	rows := tx.QueryRow(` 
		SELECT id, title, author, forum, message, votes, slug, created
		FROM threads
		WHERE id = $1
	`, id)

	var thread models.Thread
	err := rows.Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &thread.Created)
	if err != nil {
		return nil, errors.ThreadNotFound
	}

	database.CommitTransaction(tx)
	return &thread, nil
}
