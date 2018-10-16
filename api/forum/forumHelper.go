package forum

// func CreateOrGetExistingForum(forum *models.Forum) (*models.Forum, error) {
// 	// tx := DBconn.StartTransation()
// 	// defer tx.Rollback()

// 	// resultRows := tx.QueryRow(
// 	// 	`
// 	// 	INSERT
// 	// 	INTO forums ("slug", "title", "user)
// 	// 	VALUES ($1, $2, (
// 	// 			SELECT nickname FROM users WHERE nickname = $3
// 	// 		)
// 	// 	RETURNING "user
// 	// 	)
// 	// 	`,
// 	// 	forum.Slug, forum.Title, forum.User)

// }
