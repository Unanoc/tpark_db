package helpers

import (
	"bytes"
	"fmt"
	"strconv"
	"time"
	"tpark_db/database"
	"tpark_db/errors"
	"tpark_db/models"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/pgtype"
)

// ThreadCreateHelper inserts posts into table Posts.
func ThreadCreateHelper(posts *models.Posts, slugOrID string) (*models.Posts, error) {
	threadByID, err := GetThreadBySlugOrIDHelper(slugOrID)
	if err != nil {
		return nil, err
	}

	if len(*posts) == 0 {
		return posts, nil
	}

	created := time.Now().Format("2006-01-02 15:04:05")

	var sqlInsertPosts = "INSERT INTO posts (author, created, message, thread, parent, forum, path) VALUES "
	valueTemplate := "('%s', '%s', '%s', %d, %d, '%s', (SELECT path FROM posts WHERE id = %d) || (select currval(pg_get_serial_sequence('posts', 'id'))))%s"

	for i, post := range *posts {
		if authorExists(post.Author) {
			return nil, errors.UserNotFound
		}
		if parentExitsInOtherThread(post.Parent, threadByID.Id) || parentNotExists(post.Parent) {
			return nil, errors.PostParentNotFound
		}

		if i < len(*posts)-1 {
			sqlInsertPosts += fmt.Sprintf(valueTemplate, post.Author, created, post.Message, threadByID.Id, post.Parent, threadByID.Forum, post.Parent, ",")
		} else {
			sqlInsertPosts += fmt.Sprintf(valueTemplate, post.Author, created, post.Message, threadByID.Id, post.Parent, threadByID.Forum, post.Parent, "")
		}
	}

	queryRows, _ := database.DB.Conn.Query(sqlInsertPosts + "RETURNING author, created, forum, id, message, parent, thread")
	if err != nil {
		return nil, err
	}

	insertedPosts := models.Posts{}
	for queryRows.Next() {
		insertedPost := models.Post{}
		_ = queryRows.Scan(
			&insertedPost.Author,
			&insertedPost.Created,
			&insertedPost.Forum,
			&insertedPost.Id,
			&insertedPost.Message,
			&insertedPost.Parent,
			&insertedPost.Thread,
		)

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
func ThreadVoteHelper(vote *models.Vote, slugOrID string) *models.Thread {
	var err error
	prevVoice := &pgtype.Int4{}
	threadID := &pgtype.Int4{}
	threadVotes := &pgtype.Int4{}
	userNickname := &pgtype.Varchar{}

	if IsNumber(slugOrID) {
		id, _ := strconv.Atoi(slugOrID)
		err = database.DB.Conn.QueryRow(sqlSelectThreadAndVoteByID, id, vote.Nickname).Scan(prevVoice, threadID, threadVotes, userNickname)
	} else {
		err = database.DB.Conn.QueryRow(sqlSelectThreadAndVoteBySlug, slugOrID, vote.Nickname).Scan(prevVoice, threadID, threadVotes, userNickname)
	}
	if err != nil {
		return nil
	}
	if threadID.Status != pgtype.Present || userNickname.Status != pgtype.Present {
		return nil
	}
	var prevVoiceInt int32
	if prevVoice.Status == pgtype.Present {
		prevVoiceInt = int32(prevVoice.Int)
		_, err = database.DB.Conn.Exec(sqlUpdateVote, threadID.Int, userNickname.String, vote.Voice)
	} else {
		_, err = database.DB.Conn.Exec(sqlInsertVote, threadID.Int, userNickname.String, vote.Voice)
	}
	newVotes := threadVotes.Int + (int32(vote.Voice) - prevVoiceInt)
	if err != nil {
		return nil
	}
	thread := &models.Thread{}
	slugNullable := &pgtype.Varchar{}
	err = database.DB.Conn.QueryRow(sqlUpdateThreadVotes, newVotes, threadID.Int).Scan(&thread.Author, &thread.Created, &thread.Forum, &thread.Message, slugNullable, &thread.Title, &thread.Id, &thread.Votes)
	thread.Slug = slugNullable.String
	if err != nil {
		return nil
	}

	return thread
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
