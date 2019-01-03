package helpers

// UserCreate
const sqlInsertUser = `
	INSERT
	INTO users ("nickname", "fullname", "about", "email")
	VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING
`
const sqlSelectUserByEmailOrNickName = `
	SELECT "nickname", "fullname", "about", "email"
	FROM users
	WHERE "email" = $1 OR "nickname" = $2
`

// UserUpdate
const sqlUpdateUser = `
	UPDATE users
	SET fullname = coalesce(nullif($2, ''), fullname),
		about    = coalesce(nullif($3, ''), about),
		email    = coalesce(nullif($4, ''), email)
	WHERE "nickname" = $1
	RETURNING fullname, about, email, nickname
`

// UserGetOne
const sqlSelectUserByNickname = `
	SELECT "nickname", "fullname", "about", "email"
	FROM users
	WHERE "nickname" = $1
`
