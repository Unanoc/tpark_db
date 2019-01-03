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

func init() {
	ForumsCount = new(int32)
	PostsCount = new(int32)
	ThreadsCount = new(int32)
	UsersCount = new(int32)
}

func resetTablesCount() {
	atomic.StoreInt32(ThreadsCount, 0)
	atomic.StoreInt32(PostsCount, 0)
	atomic.StoreInt32(ForumsCount, 0)
	atomic.StoreInt32(UsersCount, 0)
}

// ClearHelper erases all tables.
func ClearHelper() {
	database.DB.Conn.Exec(sqlTruncateTables)
	resetTablesCount()
}

// StatusHelper returns rows counts of tables.
func StatusHelper() *models.Status {
	database.DB.Conn.QueryRow(sqlSelectCountOfTables).Scan(UsersCount, ThreadsCount, ForumsCount, PostsCount)
	currentStatus := &models.Status{
		Thread: atomic.LoadInt32(ThreadsCount),
		Post:   atomic.LoadInt32(PostsCount),
		Forum:  atomic.LoadInt32(ForumsCount),
		User:   atomic.LoadInt32(UsersCount),
	}
	return currentStatus
}
