package helpers

// ThreadCreate
const sqlInsertPost = `
	INSERT
	INTO posts (author, created, message, thread, parent, forum, path)
	VALUES ($1, $2, $3, $4, $5, $6, (SELECT path FROM posts WHERE id = $5) || (select currval(pg_get_serial_sequence('posts', 'id'))) )
	RETURNING author, created, forum, id, message, parent, thread
`

// ThreadUpdate
const sqlUpdateThread = `
	UPDATE threads
	SET title = coalesce(nullif($2, ''), title),
		message = coalesce(nullif($3, ''), message)
	WHERE slug = $1
	RETURNING id, title, author, forum, message, votes, slug, created
`

// ThreadVote
const sqlInsertVote = `
	INSERT INTO votes (thread, nickname, voice) 
	VALUES ($1, $2, $3)
`
const sqlUpdateVote = `
	UPDATE votes SET 
	voice = $3
	WHERE thread = $1 
	AND nickname = $2
`
const sqlUpdateThreadWithVote = `
	UPDATE threads SET
	votes = $1
	WHERE id = $2
	RETURNING author, created, forum, "message" , slug, title, id, votes
`

// GetThreadBySlugOrID
const sqlSelectThreadByID = `
	SELECT id, title, author, forum, message, votes, slug, created
	FROM threads
	WHERE id = $1
`
const sqlSelectThreadBySlug = `
	SELECT id, title, author, forum, message, votes, slug, created
	FROM threads
	WHERE slug = $1
`

// ThreadGetPosts
const sqlSelectPostsSinceDescLimitTree = `
	SELECT id, author, parent, message, forum, thread, created
	FROM posts
	WHERE thread = $1 AND (path < (SELECT path FROM posts WHERE id = $2::TEXT::INTEGER))
	ORDER BY path DESC
	LIMIT $3::TEXT::INTEGER
`
const sqlSelectPostsSinceDescLimitParentTree = `
	SELECT id, author, parent, message, forum, thread, created
	FROM posts
	WHERE path[1] IN (
		SELECT id
		FROM posts
		WHERE thread = $1 AND parent = 0 AND id < (SELECT path[1] FROM posts WHERE id = $2::TEXT::INTEGER)
		ORDER BY id DESC
		LIMIT $3::TEXT::INTEGER
	)
	ORDER BY path
`
const sqlSelectPostsSinceDescLimitFlat = `
	SELECT id, author, parent, message, forum, thread, created
	FROM posts
	WHERE thread = $1 AND id < $2::TEXT::INTEGER
	ORDER BY id DESC
	LIMIT $3::TEXT::INTEGER
`
const sqlSelectPostsSinceAscLimitTree = `
	SELECT id, author, parent, message, forum, thread, created
	FROM posts
	WHERE thread = $1 AND (path > (SELECT path FROM posts WHERE id = $2::TEXT::INTEGER))
	ORDER BY path
	LIMIT $3::TEXT::INTEGER
`
const sqlSelectPostsSinceAscLimitParentTree = `
	SELECT id, author, parent, message, forum, thread, created
	FROM posts
	WHERE path[1] IN (
		SELECT id
		FROM posts
		WHERE thread = $1 AND parent = 0 AND id > (SELECT path[1] FROM posts WHERE id = $2::TEXT::INTEGER)
		ORDER BY id LIMIT $3::TEXT::INTEGER
	)
	ORDER BY path
`
const sqlSelectPostsSinceAscLimitFlat = `
	SELECT id, author, parent, message, forum, thread, created
	FROM posts
	WHERE thread = $1 AND id > $2::TEXT::INTEGER
	ORDER BY id
	LIMIT $3::TEXT::INTEGER
`

const sqlSelectPostsDescLimitTree = `
	SELECT id, author, parent, message, forum, thread, created
	FROM posts
	WHERE thread = $1 
	ORDER BY path DESC
	LIMIT $2::TEXT::INTEGER
`
const sqlSelectPostsDescLimitParentTree = `
	SELECT id, author, parent, message, forum, thread, created
	FROM posts
	WHERE thread = $1 AND path[1] IN (
		SELECT path[1]
		FROM posts
		WHERE thread = $1
		GROUP BY path[1]
		ORDER BY path[1] DESC
		LIMIT $2::TEXT::INTEGER
	)
	ORDER BY path[1] DESC, path
`
const sqlSelectPostsDescLimitFlat = `
	SELECT id, author, parent, message, forum, thread, created
	FROM posts
	WHERE thread = $1
	ORDER BY id DESC
	LIMIT $2::TEXT::INTEGER
`
const sqlSelectPostsAscLimitTree = `
	SELECT id, author, parent, message, forum, thread, created
	FROM posts
	WHERE thread = $1 
	ORDER BY path
	LIMIT $2::TEXT::INTEGER
`
const sqlSelectPostsAscLimitParentTree = `
	SELECT id, author, parent, message, forum, thread, created
	FROM posts
	WHERE thread = $1 AND path[1] IN (
		SELECT path[1] 
		FROM posts 
		WHERE thread = $1 
		GROUP BY path[1]
		ORDER BY path[1]
		LIMIT $2::TEXT::INTEGER
	)
	ORDER BY path
`
const sqlSelectPostsAscLimitFlat = `
	SELECT id, author, parent, message, forum, thread, created
	FROM posts
	WHERE thread = $1 
	ORDER BY id
	LIMIT $2::TEXT::INTEGER
`
