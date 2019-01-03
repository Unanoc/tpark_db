package helpers

import (
	"bytes"
	"strconv"
	"time"
	"tpark_db/database"
	"tpark_db/errors"
	"tpark_db/models"

	"github.com/jackc/pgx"
)

// ThreadCreateHelper inserts thread into table THREADS.
func ThreadCreateHelper(posts *models.Posts, slugOrID string) (*models.Posts, error) {
	threadByID, err := GetThreadBySlugOrIDHelper(slugOrID)
	if err != nil {
		return nil, err
	}

	created := time.Now()
	insertedPosts := models.Posts{}
	for _, post := range *posts {
		if AuthorExists(post.Author) {
			return nil, errors.UserNotFound
		}

		if parentExitsInOtherThread(post.Parent, threadByID.Id) || parentNotExists(post.Parent) {
			return nil, errors.PostParentNotFound
		}

		insertedPost := models.Post{}
		err := database.DB.Conn.QueryRow(sqlInsertPost,
			post.Author,
			created,
			post.Message,
			threadByID.Id,
			post.Parent,
			threadByID.Forum,
		).Scan(
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

	return &insertedPosts, nil
}

// ThreadUpdateHelper updates thread.
func ThreadUpdateHelper(thread *models.ThreadUpdate, slugOrID string) (*models.Thread, error) {
	threadFound, err := GetThreadBySlugOrIDHelper(slugOrID)
	if err != nil {
		return nil, err
	}

	updatedThread := models.Thread{}

	err = database.DB.Conn.QueryRow(sqlUpdateThread,
		&threadFound.Slug,
		&thread.Title,
		&thread.Message,
	).Scan(
		&updatedThread.Id,
		&updatedThread.Title,
		&updatedThread.Author,
		&updatedThread.Forum,
		&updatedThread.Message,
		&updatedThread.Votes,
		&updatedThread.Slug,
		&updatedThread.Created,
	)

	if err != nil {
		return nil, err
	}

	return &updatedThread, nil
}

// ThreadVoteHelper inserts votes into table VOTES.
func ThreadVoteHelper(v *models.Vote, slugOrID string) (*models.Thread, error) {
	foundVote, _ := CheckThreadVotesByNickname(v.Nickname)
	thread, err := GetThreadBySlugOrIDHelper(slugOrID)
	if err != nil {
		return nil, err
	}

	var editedThread models.Thread
	var threadVoices int

	if foundVote == nil {
		if _, err = database.DB.Conn.Exec(sqlInsertVote, &v.Nickname, &v.Voice); err != nil {
			return nil, errors.ThreadNotFound
		}
		threadVoices = thread.Votes + v.Voice // counting of votes
	} else {
		if _, err = database.DB.Conn.Exec(sqlUpdateVote, &v.Nickname, &v.Voice); err != nil {
			return nil, err
		}
		threadVoices = thread.Votes + v.Voice - foundVote.Voice // recounting of votes with old voice
	}

	err = database.DB.Conn.QueryRow(sqlUpdateThreadWithVote, &threadVoices, &thread.Slug).Scan(
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

	return &editedThread, nil
}

// GetThreadBySlugOrIDHelper selects thread by id.
func GetThreadBySlugOrIDHelper(slugOrID string) (*models.Thread, error) {
	var err error
	var thread models.Thread
	var rows *pgx.Row

	if IsNumber(slugOrID) {
		id, _ := strconv.Atoi(slugOrID)
		rows = database.DB.Conn.QueryRow(sqlSelectThreadByID, id)
	} else {
		rows = database.DB.Conn.QueryRow(sqlSelectThreadBySlug, slugOrID)
	}

	err = rows.Scan(
		&thread.Id,
		&thread.Title,
		&thread.Author,
		&thread.Forum,
		&thread.Message,
		&thread.Votes,
		&thread.Slug,
		&thread.Created,
	)
	if err != nil {
		return nil, errors.ThreadNotFound
	}

	return &thread, nil
}

// ThreadGetPostsHelper selects posts from table POSTS with filters.
func ThreadGetPostsHelper(slugOrID string, limit, since, sort, desc []byte) (*models.Posts, error) {
	thread, err := GetThreadBySlugOrIDHelper(slugOrID)
	if err != nil {
		return nil, err
	}
	var queryRows *pgx.Rows

	if since != nil {
		if bytes.Equal([]byte("true"), desc) {
			switch string(sort) {
			case "tree":
				queryRows, err = database.DB.Conn.Query(sqlSelectPostsSinceDescLimitTree, thread.Id, since, limit)
			case "parent_tree":
				queryRows, err = database.DB.Conn.Query(sqlSelectPostsSinceDescLimitParentTree, thread.Id, since, limit)
			default:
				queryRows, err = database.DB.Conn.Query(sqlSelectPostsSinceDescLimitFlat, thread.Id, since, limit)
			}
		} else {
			switch string(sort) {
			case "tree":
				queryRows, err = database.DB.Conn.Query(sqlSelectPostsSinceAscLimitTree, thread.Id, since, limit)
			case "parent_tree":
				queryRows, err = database.DB.Conn.Query(sqlSelectPostsSinceAscLimitParentTree, thread.Id, since, limit)
			default:
				queryRows, err = database.DB.Conn.Query(sqlSelectPostsSinceAscLimitFlat, thread.Id, since, limit)
			}
		}
	} else {
		if bytes.Equal([]byte("true"), desc) {
			switch string(sort) {
			case "tree":
				queryRows, err = database.DB.Conn.Query(sqlSelectPostsDescLimitTree, thread.Id, limit)
			case "parent_tree":
				queryRows, err = database.DB.Conn.Query(sqlSelectPostsDescLimitParentTree, thread.Id, limit)
			default:
				queryRows, err = database.DB.Conn.Query(sqlSelectPostsDescLimitFlat, thread.Id, limit)
			}
		} else {
			switch string(sort) {
			case "tree":
				queryRows, err = database.DB.Conn.Query(sqlSelectPostsAscLimitTree, thread.Id, limit)
			case "parent_tree":
				queryRows, err = database.DB.Conn.Query(sqlSelectPostsAscLimitParentTree, thread.Id, limit)
			default:
				queryRows, err = database.DB.Conn.Query(sqlSelectPostsAscLimitFlat, thread.Id, limit)
			}
		}
	}
	defer queryRows.Close()

	if err != nil {
		return nil, err
	}

	posts := models.Posts{}
	for queryRows.Next() {
		post := models.Post{}

		if err = queryRows.Scan(
			&post.Id,
			&post.Author,
			&post.Parent,
			&post.Message,
			&post.Forum,
			&post.Thread,
			&post.Created,
		); err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}

	return &posts, nil
}
