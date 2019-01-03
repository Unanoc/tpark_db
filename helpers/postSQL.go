package helpers

// PostUpdate
const sqlUpdatePost = `
	UPDATE posts 
	SET message = COALESCE($2, message), "isEdited" = ($2 IS NOT NULL AND $2 <> message) 
	WHERE id = $1 
	RETURNING author::text, created, forum, "isEdited", thread, message
`

// PostGetOneByID
const sqlSelectByID = `
	SELECT id, author, message, forum, thread, created, "isEdited" 
	FROM posts 
	WHERE id = $1
`
