package helpers

// Clear
const sqlTruncateTables = `
	TRUNCATE users, forums, threads, posts, votes;
`

// Status
const sqlSelectCountOfTables = `
	SELECT *
	FROM (SELECT COUNT(*) FROM "users") as "users"
	CROSS JOIN (SELECT COUNT(*) FROM "threads") as threads
	CROSS JOIN (SELECT COUNT(*) FROM "forums") as forums
	CROSS JOIN (SELECT COUNT(*) FROM "posts") as posts
`
