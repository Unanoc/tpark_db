package helpers

import (
	"sync/atomic"
	"tpark_db/database"
	"tpark_db/models"
)

// atomic
var (
	ForumsCount  *int32
	PostsCount   *int32
	ThreadsCount *int32
	UsersCount   *int32
)

const (
	sqlClearTables = `
	TRUNCATE users, forums, threads, posts, votes, userforum;`
	sqlCountOfTables = `
	SELECT *
	FROM (SELECT COUNT(*) FROM "users") as "users"
	CROSS JOIN (SELECT COUNT(*) FROM "threads") as threads
	CROSS JOIN (SELECT COUNT(*) FROM "forums") as forums
	CROSS JOIN (SELECT COUNT(*) FROM "posts") as posts`
)

func init() {
	ForumsCount = new(int32)
	PostsCount = new(int32)
	ThreadsCount = new(int32)
	UsersCount = new(int32)
	initStatus()
}

func initStatus() {
	database.DB.Conn.QueryRow(sqlCountOfTables).Scan(UsersCount, ThreadsCount, ForumsCount, PostsCount)
}

func resetTablesCount() {
	atomic.StoreInt32(ThreadsCount, 0)
	atomic.StoreInt32(PostsCount, 0)
	atomic.StoreInt32(ForumsCount, 0)
	atomic.StoreInt32(UsersCount, 0)
}

// ClearHelper erases all tables.
func ClearHelper() {
	database.DB.Conn.Exec(sqlClearTables)
	resetTablesCount()
}

// StatusHelper returns rows counts of tables;
func StatusHelper() *models.Status {
	currentStatus := &models.Status{
		Thread: atomic.LoadInt32(ThreadsCount),
		Post:   atomic.LoadInt32(PostsCount),
		Forum:  atomic.LoadInt32(ForumsCount),
		User:   atomic.LoadInt32(UsersCount),
	}
	return currentStatus
}
