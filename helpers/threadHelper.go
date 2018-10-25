package helpers

import (
	"time"
	"tpark_db/database"
	"tpark_db/errors"
	"tpark_db/models"

	"github.com/jackc/pgx"
)

func ThreadCreateHelper(posts *models.Posts, slugOrId string) (*models.Posts, error) {
	if len(*posts) == 0 {
		return nil, errors.NoPostsForCreate
	}

	threadByID, err := GetThreadBySlugOrId(slugOrId)
	if err != nil {
		return nil, err
	}

	tx := database.StartTransaction()
	defer tx.Rollback()

	created := time.Now()
	insertedPosts := models.Posts{}
	for _, post := range *posts {
		var rows *pgx.Row

		if post.Parent != 0 {
			// TODO
		} else {
			rows = tx.QueryRow(`
				INSERT
				INTO posts (author, created, message, thread, parent, forum)
				VALUES ($1, $2, $3, $4, $5, $6) 
				RETURNING author, created, forum, id, message, parent, thread`,
				post.Author,
				created,
				post.Message,
				threadByID.Id,
				post.Parent,
				threadByID.Forum,
			)

			insertedPost := models.Post{}
			err := rows.Scan(
				&insertedPost.Author,
				&insertedPost.Created,
				&insertedPost.Forum,
				&insertedPost.Id,
				&insertedPost.Message,
				&insertedPost.Parent,
				&insertedPost.Thread,
			)
			if err != nil {
				return nil, err
			}
			insertedPosts = append(insertedPosts, &insertedPost)
		}
	}

	database.CommitTransaction(tx)
	return &insertedPosts, nil
}

func ThreadVoteHelper(v *models.Vote, slugOrId string) (*models.Thread, error) {
	tx := database.StartTransaction()
	defer tx.Rollback()

	foundVote, _ := CheckThreadVotesByNickname(v.Nickname)
	thread, err := GetThreadBySlugOrId(slugOrId)

	if err != nil {
		return nil, err
	}

	var rows *pgx.Row
	var editedThread models.Thread
	var threadVoices int

	if foundVote == nil { // like by user did not exist
		_, err := tx.Exec(`
			INSERT INTO votes (nickname, voice) VALUES ($1, $2)`,
			&v.Nickname, &v.Voice)

		if err != nil {
			return nil, err
		}

		threadVoices = thread.Votes + v.Voice // counting of votes

		rows = tx.QueryRow(`
			UPDATE threads
			SET votes = $1
			WHERE slug = $2
			RETURNING id, title, author, forum, message, votes, slug, created`, &threadVoices, &thread.Slug,
		)

		err = rows.Scan(
			&editedThread.Id,
			&editedThread.Title,
			&editedThread.Author,
			&editedThread.Forum,
			&editedThread.Message,
			&editedThread.Votes,
			&editedThread.Slug,
			&editedThread.Created,
		)

		if err != nil {
			return nil, err
		}

	} else {
		oldVote, _ := CheckThreadVotesByNickname(v.Nickname)

		if _, err := tx.Exec(`
			UPDATE votes 
			SET voice = $2
			WHERE nickname = $1`,
			&v.Nickname, &v.Voice); err != nil {
			return nil, err
		}

		threadVoices = thread.Votes + v.Voice - oldVote.Voice // counting of votes

		rows = tx.QueryRow(`
			UPDATE threads
			SET votes = $1
			WHERE slug = $2
			RETURNING id, title, author, forum, message, votes, slug, created`, &threadVoices, &thread.Slug,
		)

		err = rows.Scan(
			&editedThread.Id,
			&editedThread.Title,
			&editedThread.Author,
			&editedThread.Forum,
			&editedThread.Message,
			&editedThread.Votes,
			&editedThread.Slug,
			&editedThread.Created,
		)

		if err != nil {
			return nil, err
		}
	}

	database.CommitTransaction(tx)
	return &editedThread, nil
}
