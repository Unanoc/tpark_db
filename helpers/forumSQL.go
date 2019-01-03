package helpers

// ForumCreate
const sqlInsertForum = `
	INSERT
	INTO forums (slug, title, "user")
	VALUES ($1, $2, (SELECT nickname FROM users WHERE nickname = $3)) 
	RETURNING "user"
`

// ForumCreateThread
const sqlInsertThread = `
	INSERT
	INTO threads (author, created, message, title, slug, forum)
	VALUES ($1, $2, $3, $4, $5, (SELECT slug FROM forums WHERE slug = $6)) 
	RETURNING author, created, forum, id, message, title
`

// ForumGetBySlug
const sqlSelectForumBySlug = `
	SELECT slug, title, "user", 
		(SELECT COUNT(*) FROM posts WHERE forum = $1), 
		(SELECT COUNT(*) FROM threads WHERE forum = $1)
	FROM forums
	WHERE slug = $1
`

// ForumGetThreads
const sqlSelectThreadsSinceDescLimit = `
	SELECT author, created, forum, id, message, slug, title, votes
	FROM threads
	WHERE forum = $1 AND created <= $2::TEXT::TIMESTAMPTZ
	ORDER BY created DESC
	LIMIT $3::TEXT::INTEGER
`
const sqlSelectThreadsSinceAscLimit = `
	SELECT author, created, forum, id, message, slug, title, votes
	FROM threads
	WHERE forum = $1 AND created >= $2::TEXT::TIMESTAMPTZ
	ORDER BY created
	LIMIT $3::TEXT::INTEGER
`
const sqlSelectThreadDescLimit = `
	SELECT author, created, forum, id, message, slug, title, votes
	FROM threads
	WHERE forum = $1
	ORDER BY created DESC
	LIMIT $2::TEXT::INTEGER
`
const sqlSelectThreadAscLimit = `
	SELECT author, created, forum, id, message, slug, title, votes
	FROM threads
	WHERE forum = $1
	ORDER BY created
	LIMIT $2::TEXT::INTEGER
`

// ForumGetUsers
const sqlSelectUsersDescSinceLimit = `
	SELECT nickname, fullname, about, email
	FROM users
	WHERE nickname IN (
			SELECT author FROM threads WHERE forum = $1
			UNION
			SELECT author FROM posts WHERE forum = $1
		) 
		AND LOWER(nickname) < LOWER($2::TEXT)
	ORDER BY nickname DESC
	LIMIT $3::TEXT::INTEGER
`
const sqlSelectUsersAscSinceLimit = `
	SELECT nickname, fullname, about, email
	FROM users
	WHERE nickname IN (
			SELECT author FROM threads WHERE forum = $1
			UNION
			SELECT author FROM posts WHERE forum = $1
		)  
		AND LOWER(nickname) > LOWER($2::TEXT)
	ORDER BY nickname
	LIMIT $3::TEXT::INTEGER
`
const sqlSelectUsersDescLimit = `
	SELECT nickname, fullname, about, email
	FROM users
	WHERE nickname IN (
			SELECT author FROM threads WHERE forum = $1
			UNION
			SELECT author FROM posts WHERE forum = $1
		) 
	ORDER BY nickname DESC
	LIMIT $2::TEXT::INTEGER
`
const sqlSelectUsersAscLimit = `
	SELECT nickname, fullname, about, email
	FROM users
	WHERE nickname IN (
			SELECT author FROM threads WHERE forum = $1
			UNION
			SELECT author FROM posts WHERE forum = $1
		) 
	ORDER BY nickname
	LIMIT $2::TEXT::INTEGER
`
