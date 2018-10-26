package helpers

import (
	"log"
	"sync"
	"tpark_db/database"
	"tpark_db/models"
)

func getCountOfTable(wg *sync.WaitGroup, table string, status *models.Status) {
	defer wg.Done()
	queryString := "SELECT COUNT(*) FROM " + table
	var count int
	tx := database.StartTransaction()
	defer tx.Rollback()
	rows := tx.QueryRow(queryString)
	_ = rows.Scan(&count)
	database.CommitTransaction(tx)
	switch table {
	case "users":
		status.User = count
	case "forums":
		status.Forum = count
	case "threads":
		status.Thread = count
	case "posts":
		status.Post = count
	default:
		log.Panic("wrong type of table")
	}
}

func StatusHelper() (*models.Status, error) {
	status := models.Status{}
	listTables := []string{"users", "forums", "threads", "posts"}
	wg := &sync.WaitGroup{}
	for _, table := range listTables {
		wg.Add(1)
		go getCountOfTable(wg, table, &status)
	}
	wg.Wait()
	return &status, nil
}
